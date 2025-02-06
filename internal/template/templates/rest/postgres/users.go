package postgres

import (
	"database/sql"
	"log/slog"

	"example.com/rest"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) rest.UserService {
	return &userService{db: db}
}

func (s *userService) Create(req *rest.UserCredentials) (*rest.User, error) {

	fields := req.Validate()
	if fields != nil {
		return nil, rest.ErrorfWithFields(rest.INVALID_ERROR, "invalid input", fields)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, rest.Errorf(rest.INTERNAL_ERROR, "internal error")
	}

	var user rest.User
	err = s.db.QueryRow(`
        INSERT INTO users (username, password_hash)
        VALUES ($1, $2)
        RETURNING id, username, created_at, updated_at`,
		req.Username, string(hashedPassword),
	).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		// Check for unique constraint violation, i.e. username already exists
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, rest.Errorf(rest.CONFLICT_ERROR, "username already exists")
		}
		slog.Error("failed to create user", "error", err)
		return nil, rest.Errorf(rest.INTERNAL_ERROR, "internal error")
	}

	return &user, nil
}

func (s *userService) Get(id int) (*rest.User, error) {
	var user rest.User
	err := s.db.QueryRow(`
        SELECT id, username, created_at 
        FROM users 
        WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, rest.Errorf(rest.NOTFOUND_ERROR, "user not found")
		}
		slog.Error("failed to get user", "error", err)
		return nil, rest.Errorf(rest.INTERNAL_ERROR, "internal error")
	}

	return &user, nil
}

func (s *userService) GetByUsername(username string) (*rest.User, error) {
	var user rest.User
	err := s.db.QueryRow(`
        SELECT id, username, created_at 
        FROM users 
        WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, rest.Errorf(rest.NOTFOUND_ERROR, "user not found")
		}
		slog.Error("failed to get user", "error", err)
		return nil, rest.Errorf(rest.INTERNAL_ERROR, "internal error")
	}

	return &user, nil
}

func (s *userService) Authenticate(req *rest.UserCredentials) (*rest.User, error) {

	fields := req.Validate()
	if fields != nil {
		return nil, rest.ErrorfWithFields(rest.INVALID_ERROR, "invalid input", fields)
	}

	var user rest.User
	err := s.db.QueryRow(`
        SELECT id, username, password_hash, created_at 
        FROM users 
        WHERE username = $1`,
		req.Username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, rest.Errorf(rest.UNAUTHORIZED_ERROR, "invalid credentials")
		}
		slog.Error("failed to authenticate user", "error", err)
		return nil, rest.Errorf(rest.INTERNAL_ERROR, "internal error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, rest.Errorf(rest.UNAUTHORIZED_ERROR, "invalid credentials")
	}

	return &user, nil
}

func (s *userService) Update(id int, req *rest.UpdateUser) (*rest.User, error) {

	fields := req.Validate()
	if fields != nil {
		return nil, rest.ErrorfWithFields(rest.INVALID_ERROR, "invalid input", fields)
	}

	var user rest.User
	err := s.db.QueryRow(`
        UPDATE users 
        SET username = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
        RETURNING id, username, created_at, updated_at`,
		req.Username, id,
	).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, rest.Errorf(rest.NOTFOUND_ERROR, "user not found")
		}
		slog.Error("failed to update user", "error", err)
		return nil, rest.Errorf(rest.INTERNAL_ERROR, "internal error")
	}

	return &user, nil
}

func (s *userService) Delete(id int) error {
	result, err := s.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		slog.Error("error deleting user", "error", err)
		return rest.Errorf(rest.INTERNAL_ERROR, "internal error")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		slog.Error("error checking rows affected", "error", err)
		return rest.Errorf(rest.INTERNAL_ERROR, "internal error")
	}

	if rows == 0 {
		return rest.Errorf(rest.NOTFOUND_ERROR, "user not found")
	}

	return nil
}
