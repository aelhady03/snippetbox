package models

import (
	"database/sql"
	"errors"
	"time"
)

// Snippet represents an individual snippet.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel type wraps a DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

func NewSnippetModel(db *sql.DB) *SnippetModel {
	return &SnippetModel{
		DB: db,
	}
}

// Insert insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get gets an existing snippet from the database.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
	FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	var snippet Snippet

	err := m.DB.QueryRow(stmt, id).Scan(
		&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		}
		return Snippet{}, err
	}

	return snippet, nil
}

// Latest returns the latest 10 snippets from the database.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
	FROM snippets WHERE expires > UTC_TIMESTAMP()
	ORDER BY id DESC LIMIT 10`

	var snippets []Snippet

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var snippet Snippet

		err := rows.Scan(
			&snippet.ID,
			&snippet.Title,
			&snippet.Content,
			&snippet.Created,
			&snippet.Expires,
		)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
