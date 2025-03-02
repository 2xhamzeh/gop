package postgres

import (
	"database/sql"
	"log/slog"

	"example.com/rest/internal/domain"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewUserService(db *sql.DB, logger *slog.Logger) *userService {
	return &userService{
		db:     db,
		logger: logger,
	}
}

func (s *userService) Create(req *domain.UserCredentials) (*domain.User, error) {

	fields := req.Validate()
	if fields != nil {
		return nil, domain.ErrorfWithFields(domain.INVALID_ERROR, "invalid input", fields)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}

	var user domain.User
	err = s.db.QueryRow(`
        INSERT INTO users (username, password_hash)
        VALUES ($1, $2)
        RETURNING id, username, created_at, updated_at`,
		req.Username, string(hashedPassword),
	).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		// Check for unique constraint violation, i.e. username already exists
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, domain.Errorf(domain.CONFLICT_ERROR, "username already exists")
		}
		s.logger.Error("failed to create user", "error", err)
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}

	return &user, nil
}

func (s *userService) Get(id int) (*domain.User, error) {
	var user domain.User
	err := s.db.QueryRow(`
        SELECT id, username, created_at 
        FROM users 
        WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.Errorf(domain.NOTFOUND_ERROR, "user not found")
		}
		s.logger.Error("failed to get user", "error", err)
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}

	return &user, nil
}

func (s *userService) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := s.db.QueryRow(`
        SELECT id, username, created_at 
        FROM users 
        WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.Errorf(domain.NOTFOUND_ERROR, "user not found")
		}
		s.logger.Error("failed to get user", "error", err)
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}

	return &user, nil
}

func (s *userService) Authenticate(req *domain.UserCredentials) (*domain.User, error) {

	fields := req.Validate()
	if fields != nil {
		return nil, domain.ErrorfWithFields(domain.INVALID_ERROR, "invalid input", fields)
	}

	var user domain.User
	err := s.db.QueryRow(`
        SELECT id, username, password_hash, created_at 
        FROM users 
        WHERE username = $1`,
		req.Username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid credentials")
		}
		s.logger.Error("failed to authenticate user", "error", err)
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid credentials")
	}

	return &user, nil
}

func (s *userService) Update(id int, req *domain.UpdateUser) (*domain.User, error) {

	fields := req.Validate()
	if fields != nil {
		return nil, domain.ErrorfWithFields(domain.INVALID_ERROR, "invalid input", fields)
	}

	var user domain.User
	err := s.db.QueryRow(`
        UPDATE users 
        SET username = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
        RETURNING id, username, created_at, updated_at`,
		req.Username, id,
	).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.Errorf(domain.NOTFOUND_ERROR, "user not found")
		}
		s.logger.Error("failed to update user", "error", err)
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}

	return &user, nil
}

func (s *userService) Delete(id int) error {
	result, err := s.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		s.logger.Error("error deleting user", "error", err)
		return domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("error checking rows affected", "error", err)
		return domain.Errorf(domain.INTERNAL_ERROR, "internal error")
	}

	if rows == 0 {
		return domain.Errorf(domain.NOTFOUND_ERROR, "user not found")
	}

	return nil
}
