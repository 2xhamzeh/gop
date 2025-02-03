package postgres

import (
	"database/sql"
	"log/slog"

	"example.com/app"
)

type noteService struct {
	db *sql.DB
}

func NewNoteService(db *sql.DB) app.NoteService {
	return &noteService{db: db}
}

func (s *noteService) Create(userID int, req *app.CreateNote) (*app.Note, error) {

	fields := req.Validate()
	if fields != nil {
		return nil, app.ErrorfWithFields(app.INVALID_ERROR, "invalid input", fields)
	}

	var note app.Note
	err := s.db.QueryRow(`
        INSERT INTO notes (user_id, title, content)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, title, content, created_at, updated_at`,
		userID, req.Title, req.Content,
	).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		slog.Error("failed to create note", "error", err)
		return nil, app.Errorf(app.INTERNAL_ERROR, "internal error")
	}

	return &note, nil
}

func (s *noteService) GetAll(userID int) ([]app.Note, error) {
	rows, err := s.db.Query(`
        SELECT id, title, content, created_at, updated_at
        FROM notes
        WHERE user_id = $1
        ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app.Errorf(app.NOTFOUND_ERROR, "no notes found")
		}
		slog.Error("failed to get notes", "error", err)
		return nil, app.Errorf(app.INTERNAL_ERROR, "internal error")
	}
	defer rows.Close()

	data := []app.Note{}
	for rows.Next() {
		var note app.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			slog.Error("failed to scan note", "error", err)
			return nil, app.Errorf(app.INTERNAL_ERROR, "internal error")
		}
		data = append(data, note)
	}
	if err = rows.Err(); err != nil {
		slog.Error("failed to iterate over notes", "error", err)
		return nil, app.Errorf(app.INTERNAL_ERROR, "internal error")
	}

	return data, nil
}

func (s *noteService) Update(userID int, noteID int, req *app.UpdateNote) (*app.Note, error) {

	fields := req.Validate()
	if fields != nil {
		return nil, app.ErrorfWithFields(app.INVALID_ERROR, "invalid input", fields)
	}

	var note app.Note
	err := s.db.QueryRow(`
        UPDATE notes
        SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3 AND user_id = $4
        RETURNING id, title, content, created_at, updated_at`,
		req.Title, req.Content, noteID, userID,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app.Errorf(app.NOTFOUND_ERROR, "note not found")
		}
		slog.Error("failed to update note", "error", err)
		return nil, app.Errorf(app.INTERNAL_ERROR, "internal error")
	}

	return &note, nil
}

func (s *noteService) Delete(userID int, noteID int) error {
	result, err := s.db.Exec(`
        DELETE FROM notes 
        WHERE id = $1 AND user_id = $2`,
		noteID, userID,
	)
	if err != nil {
		slog.Error("failed to delete note", "error", err)
		return app.Errorf(app.INTERNAL_ERROR, "internal error")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to check rows affected", "error", err)
		return app.Errorf(app.INTERNAL_ERROR, "internal error")
	}

	if rows == 0 {
		return app.Errorf(app.NOTFOUND_ERROR, "note not found")
	}

	return nil
}
