package mysql

import (
	"database/sql"
	"errors"
	"ntguilty.me/notes-app/pkg/models"
)

type NoteModel struct {
	DB *sql.DB
}

func (m *NoteModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO notes (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

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

func (m *NoteModel) Get(id int) (*models.Note, error) {
	stmt := `SELECT id, title, content, created, expires FROM notes 
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// This return a pointer to a sql.Row which holds the result from the database
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new empty Note struct and then
	// use row.Scan() to copy values from sql.Row to the corresponding field in the
	// Note struct
	s := &models.Note{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	// This show models.ErrNoRecord instead of sql.ErrNoRows because we want our app
	// not to be concerned about datastore-specific errors
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// This will return the 10 most recently created notes.
func (m *NoteModel) Latest() ([]*models.Note, error) {
	return nil, nil
}
