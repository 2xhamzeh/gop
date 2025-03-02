package postgres

import (
	"database/sql"
	"log/slog"

	"example.com/rest/internal/domain"
)

type noteService struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewNoteService(db *sql.DB, logger *slog.Logger) *noteService {
	return &noteService{
		db:     db,
		logger: logger,
	}
}

func (s *noteService) Create(userID int, req *domain.CreateNote) (*domain.Note, error) {

	fields := req.Validate()
	if fields != nil {
		return nil, domain.ErrorfWithFields(domain.INVALID_ERROR, "invalid input", fields)
	}

	var note domain.Note
	err := s.db.QueryRow(`
        INSERT INTO notes (user_id, title, content)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, title, content, created_at, updated_at`,
		userID, req.Title, req.Content,
	).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		s.logger.Error("failed to create note", "error", err)
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}

	return &note, nil
}

func (s *noteService) GetAll(userID int) ([]domain.Note, error) {
	rows, err := s.db.Query(`
        SELECT id, title, content, created_at, updated_at
        FROM notes
        WHERE user_id = $1
        ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.Errorf(domain.NOTFOUND_ERROR, "no notes found")
		}
		s.logger.Error("failed to get notes", "error", err)
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}
	defer rows.Close()

	data := []domain.Note{}
	for rows.Next() {
		var note domain.Note
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			s.logger.Error("failed to scan note", "error", err)
			return nil, domain.Errorf(domain.INTERNAL_ERROR, "internal error")
		}
		data = append(data, note)
	}
	if err = rows.Err(); err != nil {
		s.logger.Error("failed to iterate over notes", "error", err)
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}

	return data, nil
}

func (s *noteService) Update(userID int, noteID int, req *domain.UpdateNote) (*domain.Note, error) {

	fields := req.Validate()
	if fields != nil {
		return nil, domain.ErrorfWithFields(domain.INVALID_ERROR, "invalid input", fields)
	}

	var note domain.Note
	err := s.db.QueryRow(`
        UPDATE notes
        SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP
        WHERE id = $3 AND user_id = $4
        RETURNING id, title, content, created_at, updated_at`,
		req.Title, req.Content, noteID, userID,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.Errorf(domain.NOTFOUND_ERROR, "note not found")
		}
		s.logger.Error("failed to update note", "error", err)
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "internal error")
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
		s.logger.Error("failed to delete note", "error", err)
		return domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("failed to check rows affected", "error", err)
		return domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}

	if rows == 0 {
		return domain.Errorf(domain.NOTFOUND_ERROR, "note not found")
	}

	return nil
}
