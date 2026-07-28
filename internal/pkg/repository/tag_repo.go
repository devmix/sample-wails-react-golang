package repository

import (
	"database/sql"
	"fmt"

	"sample-wails-react-golang/internal/pkg/model"
)

// TagRepo provides data access for tag records.
type TagRepo struct {
	db *sql.DB // underlying database connection.
}

// NewTagRepo creates a new TagRepo backed by the given DB instance.
func NewTagRepo(db *DB) *TagRepo {
	return &TagRepo{db: db.conn}
}

// ListAll returns every tag ordered alphabetically by name.
func (r *TagRepo) ListAll() ([]model.Tag, error) {
	rows, err := r.db.Query(`SELECT id, name FROM tags ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("query tags: %w", err)
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
