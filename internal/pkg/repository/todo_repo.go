package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"sample-wails-react-golang/internal/pkg/model"
)

// TodoRepo provides data access for todo records.
type TodoRepo struct {
	db *sql.DB // underlying database connection.
}

// NewTodoRepo creates a new TodoRepo backed by the given DB instance.
func NewTodoRepo(db *DB) *TodoRepo {
	return &TodoRepo{db: db.conn}
}

// List returns todos with optional filtering by status, priority, or tag name.
// Results are ordered by priority (urgent first) then due date ascending.
func (r *TodoRepo) List(filter *model.TodoFilter) ([]model.Todo, error) {
	whereParts := []string{"1=1"}
	args := []interface{}{}

	if filter != nil {
		if filter.Status != nil && *filter.Status != "" {
			whereParts = append(whereParts, "t.status = ?")
			args = append(args, *filter.Status)
		}
		if filter.Priority != nil && *filter.Priority != "" {
			whereParts = append(whereParts, "t.priority = ?")
			args = append(args, *filter.Priority)
		}
		if filter.TagName != nil && *filter.TagName != "" {
			whereParts = append(whereParts, `t.id IN (SELECT tt.todo_id FROM todo_tags tt JOIN tags tg ON tt.tag_id = tg.id WHERE LOWER(tg.name) = LOWER(?))`)
			args = append(args, *filter.TagName)
		}
	}

	query := fmt.Sprintf(`
		SELECT t.id, t.title, t.description, t.status, t.priority, t.due_date, t.completed_at, t.created_at, t.updated_at,
		       tg.id AS tag_id, tg.name AS tag_name
		FROM todos t
		LEFT JOIN todo_tags tgt ON t.id = tgt.todo_id
		LEFT JOIN tags tg ON tgt.tag_id = tg.id
		WHERE %s
		ORDER BY CASE t.priority WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 END, t.due_date ASC`,
		strings.Join(whereParts, " AND "))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query todos: %w", err)
	}
	defer rows.Close()

	return r.scanTodosWithTags(rows)
}

// Get returns a single todo by ID along with its associated tags.
func (r *TodoRepo) Get(id int) (*model.Todo, error) {
	row := r.db.QueryRow(`SELECT id, title, description, status, priority, due_date, completed_at, created_at, updated_at FROM todos WHERE id = ?`, id)
	todo := &model.Todo{}
	err := row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Priority, &todo.DueDate, &todo.CompletedAt, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("todo %d not found", id)
		}
		return nil, fmt.Errorf("get todo: %w", err)
	}

	tags, err := r.getTodoTags(id)
	if err != nil {
		return nil, err
	}
	todo.Tags = tags
	return todo, nil
}

// Create inserts a new todo record within a transaction. If no priority is given,
// it defaults to "medium". Associated tags are created and linked in the same transaction.
func (r *TodoRepo) Create(input model.TodoCreateInput) (*model.Todo, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	priority := input.Priority
	if priority == "" {
		priority = "medium"
	}

	res, err := tx.Exec(`INSERT INTO todos (title, description, status, priority, due_date) VALUES (?, ?, 'pending', ?, ?)`,
		input.Title, input.Description, priority, input.DueDate)
	if err != nil {
		return nil, fmt.Errorf("insert todo: %w", err)
	}

	id, _ := res.LastInsertId()
	if err := r.setTodoTags(tx, int(id), input.TagNames); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return r.Get(int(id))
}

// Update modifies fields of an existing todo within a transaction. Only non-nil
// input fields are updated. Setting status to "completed" records the completion time.
func (r *TodoRepo) Update(id int, input model.TodoUpdateInput) (*model.Todo, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	setClauses := []string{"updated_at = strftime('%Y-%m-%dT%H:%M:%f', 'now')"}
	args := []interface{}{}

	if input.Title != nil {
		setClauses = append(setClauses, "title = ?")
		args = append(args, *input.Title)
	}
	if input.Description != nil {
		setClauses = append(setClauses, "description = ?")
		args = append(args, *input.Description)
	}
	if input.Status != nil {
		setClauses = append(setClauses, "status = ?")
		args = append(args, *input.Status)
		if *input.Status == "completed" {
			setClauses = append(setClauses, "completed_at = strftime('%Y-%m-%dT%H:%M:%f', 'now')")
		} else {
			setClauses = append(setClauses, "completed_at = NULL")
		}
	}
	if input.Priority != nil {
		setClauses = append(setClauses, "priority = ?")
		args = append(args, *input.Priority)
	}
	if input.DueDate != nil {
		setClauses = append(setClauses, "due_date = ?")
		args = append(args, *input.DueDate)
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE todos SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	if _, err := tx.Exec(query, args...); err != nil {
		return nil, fmt.Errorf("update todo: %w", err)
	}

	if len(input.TagNames) > 0 {
		if err := r.setTodoTags(tx, id, input.TagNames); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return r.Get(id)
}

// Delete removes a todo by ID. Returns an error if the todo does not exist.
func (r *TodoRepo) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("todo %d not found", id)
	}
	return nil
}

// Search finds todos whose title or description match the query using a LIKE pattern.
// Results are ordered by creation date descending.
func (r *TodoRepo) Search(query string) ([]model.Todo, error) {
	q := "%" + query + "%"
	rows, err := r.db.Query(`
		SELECT t.id, t.title, t.description, t.status, t.priority, t.due_date, t.completed_at, t.created_at, t.updated_at,
		       tg.id AS tag_id, tg.name AS tag_name
		FROM todos t
		LEFT JOIN todo_tags tgt ON t.id = tgt.todo_id
		LEFT JOIN tags tg ON tgt.tag_id = tg.id
		WHERE t.title LIKE ? OR t.description LIKE ?
		ORDER BY t.created_at DESC`, q, q)
	if err != nil {
		return nil, fmt.Errorf("search todos: %w", err)
	}
	defer rows.Close()

	return r.scanTodosWithTags(rows)
}

func (r *TodoRepo) scanTodosWithTags(rows *sql.Rows) ([]model.Todo, error) {
	todosMap := make(map[int]*model.Todo)
	var order []int

	for rows.Next() {
		var id int
		var title, description, status, priority string
		var dueDate, completedAt, createdAt, updatedAt sql.NullString
		var tagID sql.NullInt64
		var tagName sql.NullString

		err := rows.Scan(&id, &title, &description, &status, &priority, &dueDate, &completedAt, &createdAt, &updatedAt, &tagID, &tagName)
		if err != nil {
			return nil, err
		}

		todo, exists := todosMap[id]
		if !exists {
			var dueDateStr *string
			if dueDate.Valid {
				dueDateStr = &dueDate.String
			}
			var completedAtStr *string
			if completedAt.Valid {
				completedAtStr = &completedAt.String
			}
			todo = &model.Todo{
				ID:          id,
				Title:       title,
				Description: description,
				Status:      status,
				Priority:    priority,
				DueDate:     dueDateStr,
				CompletedAt: completedAtStr,
				CreatedAt:   createdAt.String,
				UpdatedAt:   updatedAt.String,
				Tags:        []model.Tag{},
			}
			todosMap[id] = todo
			order = append(order, id)
		}

		if tagID.Valid {
			todo.Tags = append(todo.Tags, model.Tag{
				ID:   int(tagID.Int64),
				Name: tagName.String,
			})
		}
	}

	result := make([]model.Todo, 0, len(order))
	for _, id := range order {
		result = append(result, *todosMap[id])
	}

	return result, rows.Err()
}

func (r *TodoRepo) getTodoTags(todoID int) ([]model.Tag, error) {
	rows, err := r.db.Query(`SELECT t.id, t.name FROM tags t JOIN todo_tags tt ON t.id = tt.tag_id WHERE tt.todo_id = ?`, todoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func (r *TodoRepo) setTodoTags(tx *sql.Tx, todoID int, names []string) error {
	if _, err := tx.Exec("DELETE FROM todo_tags WHERE todo_id = ?", todoID); err != nil {
		return fmt.Errorf("clear todo tags: %w", err)
	}

	for _, name := range names {
		if name == "" {
			continue
		}
		res, err := tx.Exec(`INSERT INTO tags (name) VALUES (?) ON CONFLICT(name) DO NOTHING`, name)
		if err != nil {
			return fmt.Errorf("insert tag: %w", err)
		}

		tagID, _ := res.LastInsertId()
		if tagID == 0 {
			var id int
			if err := tx.QueryRow(`SELECT id FROM tags WHERE LOWER(name) = LOWER(?)`, name).Scan(&id); err != nil {
				return fmt.Errorf("find existing tag: %w", err)
			}
			tagID = int64(id)
		}

		if _, err := tx.Exec(`INSERT OR IGNORE INTO todo_tags (todo_id, tag_id) VALUES (?, ?)`, todoID, tagID); err != nil {
			return fmt.Errorf("link todo tag: %w", err)
		}
	}
	return nil
}
