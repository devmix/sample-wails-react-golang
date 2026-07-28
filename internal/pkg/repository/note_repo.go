package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"sample-wails-react-golang/internal/pkg/model"
)

// NoteRepo provides data access for note records.
type NoteRepo struct {
	db *sql.DB // underlying database connection.
}

// NewNoteRepo creates a new NoteRepo backed by the given DB instance.
func NewNoteRepo(db *DB) *NoteRepo {
	return &NoteRepo{db: db.conn}
}

// List returns all notes with their tags, ordered by most recently updated.
func (r *NoteRepo) List() ([]model.Note, error) {
	rows, err := r.db.Query(`
		SELECT n.id, n.title, n.content, n.created_at, n.updated_at,
		       tg.id AS tag_id, tg.name AS tag_name
		FROM notes n
		LEFT JOIN note_tags nt ON n.id = nt.note_id
		LEFT JOIN tags tg ON nt.tag_id = tg.id
		ORDER BY n.updated_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("query notes: %w", err)
	}
	defer rows.Close()

	return r.scanNotesWithTags(rows)
}

// Get returns a single note by ID along with its associated tags.
func (r *NoteRepo) Get(id int) (*model.Note, error) {
	row := r.db.QueryRow(`SELECT id, title, content, created_at, updated_at FROM notes WHERE id = ?`, id)
	note := &model.Note{}
	err := row.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("note %d not found", id)
		}
		return nil, fmt.Errorf("get note: %w", err)
	}

	tags, err := r.getNoteTags(id)
	if err != nil {
		return nil, err
	}
	note.Tags = tags
	return note, nil
}

// Create inserts a new note within a transaction and links any provided tags.
func (r *NoteRepo) Create(input model.NoteCreateInput) (*model.Note, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO notes (title, content) VALUES (?, ?)`, input.Title, input.Content)
	if err != nil {
		return nil, fmt.Errorf("insert note: %w", err)
	}

	id, _ := res.LastInsertId()
	if err := r.setNoteTags(tx, int(id), input.TagNames); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return r.Get(int(id))
}

// Update modifies fields of an existing note within a transaction. Only non-nil
// input fields are updated. Tag associations can be replaced if provided.
func (r *NoteRepo) Update(id int, input model.NoteUpdateInput) (*model.Note, error) {
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
	if input.Content != nil {
		setClauses = append(setClauses, "content = ?")
		args = append(args, *input.Content)
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE notes SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	if _, err := tx.Exec(query, args...); err != nil {
		return nil, fmt.Errorf("update note: %w", err)
	}

	if len(input.TagNames) > 0 {
		if err := r.setNoteTags(tx, id, input.TagNames); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return r.Get(id)
}

// Delete removes a note by ID. Returns an error if the note does not exist.
func (r *NoteRepo) Delete(id int) error {
	res, err := r.db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("note %d not found", id)
	}
	return nil
}

// Search finds notes whose title or content match the query using a LIKE pattern.
// Results are ordered by update date descending.
func (r *NoteRepo) Search(query string) ([]model.Note, error) {
	q := "%" + query + "%"
	rows, err := r.db.Query(`
		SELECT n.id, n.title, n.content, n.created_at, n.updated_at,
		       tg.id AS tag_id, tg.name AS tag_name
		FROM notes n
		LEFT JOIN note_tags nt ON n.id = nt.note_id
		LEFT JOIN tags tg ON nt.tag_id = tg.id
		WHERE n.title LIKE ? OR n.content LIKE ?
		ORDER BY n.updated_at DESC`, q, q)
	if err != nil {
		return nil, fmt.Errorf("search notes: %w", err)
	}
	defer rows.Close()

	return r.scanNotesWithTags(rows)
}

func (r *NoteRepo) scanNotesWithTags(rows *sql.Rows) ([]model.Note, error) {
	notesMap := make(map[int]*model.Note)
	var order []int

	for rows.Next() {
		var id int
		var title, content, createdAt, updatedAt string
		var tagID sql.NullInt64
		var tagName sql.NullString

		err := rows.Scan(&id, &title, &content, &createdAt, &updatedAt, &tagID, &tagName)
		if err != nil {
			return nil, err
		}

		note, exists := notesMap[id]
		if !exists {
			note = &model.Note{
				ID:        id,
				Title:     title,
				Content:   content,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				Tags:      []model.Tag{},
			}
			notesMap[id] = note
			order = append(order, id)
		}

		if tagID.Valid {
			note.Tags = append(note.Tags, model.Tag{
				ID:   int(tagID.Int64),
				Name: tagName.String,
			})
		}
	}

	result := make([]model.Note, 0, len(order))
	for _, id := range order {
		result = append(result, *notesMap[id])
	}

	return result, rows.Err()
}

func (r *NoteRepo) getNoteTags(noteID int) ([]model.Tag, error) {
	rows, err := r.db.Query(`SELECT t.id, t.name FROM tags t JOIN note_tags nt ON t.id = nt.tag_id WHERE nt.note_id = ?`, noteID)
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

func (r *NoteRepo) setNoteTags(tx *sql.Tx, noteID int, names []string) error {
	if _, err := tx.Exec("DELETE FROM note_tags WHERE note_id = ?", noteID); err != nil {
		return fmt.Errorf("clear note tags: %w", err)
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

		if _, err := tx.Exec(`INSERT OR IGNORE INTO note_tags (note_id, tag_id) VALUES (?, ?)`, noteID, tagID); err != nil {
			return fmt.Errorf("link note tag: %w", err)
		}
	}
	return nil
}
