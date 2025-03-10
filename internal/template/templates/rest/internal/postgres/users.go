package postgres

import (
	"database/sql"

	"example.com/rest/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	db *sqlx.DB
}

func NewUserService(db *sqlx.DB) *userService {
	return &userService{
		db: db,
	}
}

// Create takes a UserCredentials and inserts a new user into the database.
// It returns the created user or an error if the operation fails.
func (s *userService) Create(req *domain.UserCredentials) (*domain.User, error) {
	// validate input
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "failed to hash password: %v", err)
	}

	// insert user
	var user domain.User
	err = s.db.QueryRow(`
        INSERT INTO users (email, password_hash)
        VALUES ($1, $2)
        RETURNING id, email, created_at, updated_at`,
		req.Email, string(hashedPassword),
	).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		// Check for unique constraint violation: email already exists
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, domain.Errorf(domain.CONFLICT_ERROR, "email already exists")
		}
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "failed to create user: %v", err)
	}
	return &user, nil
}

// Get takes a user ID and finds the user in the database.
// It returns the user or an error if the operation fails.
func (s *userService) Get(id int) (*domain.User, error) {
	// find user
	var user domain.User
	err := s.db.QueryRow(`
        SELECT id, email, created_at, updated_at
        FROM users 
        WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.Errorf(domain.NOTFOUND_ERROR, "user not found")
		}
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "failed to get user: %v", err)
	}

	return &user, nil
}

// Authenticate takes a UserCredentials and checks if the user exists and the password is correct.
// It returns the user or an error if the operation fails.
func (s *userService) Authenticate(req *domain.UserCredentials) (*domain.User, error) {
	// validate input
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	// find user
	var user domain.User
	err = s.db.QueryRow(`
        SELECT id, email, password_hash, created_at, updated_at
        FROM users 
        WHERE email = $1`,
		req.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid credentials")
		}
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "failed to authenticate user: %v", err)
	}

	// compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid credentials")
	}

	return &user, nil
}

// Update takes a user ID and a UserPatch and updates the user in the database.
// It returns the updated user or an error if the operation fails.
func (s *userService) Update(id int, req *domain.UserPatch) (*domain.User, error) {
	// validate input
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	// get user
	var user domain.User
	err = s.db.QueryRow(`
		SELECT email, password_hash, version
		FROM users
		WHERE id = $1`,
		id,
	).Scan(&user.Email, &user.PasswordHash, &user.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.Errorf(domain.NOTFOUND_ERROR, "user not found")
		}
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "failed to get user: %v", err)
	}

	// we could check if the new information is the same as the old one

	// compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, domain.Errorf(domain.UNAUTHORIZED_ERROR, "invalid password")
	}

	// check which fields to update
	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.NewPassword != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, domain.Errorf(domain.INTERNAL_ERROR, "failed to hash password: %v", err)
		}
		user.PasswordHash = string(hashedPassword)
	}

	// update user
	// use version to insure no concurrent updates
	err = s.db.QueryRow(`
        UPDATE users
        SET email = $1, password_hash = $2, updated_at = CURRENT_TIMESTAMP, version = version + 1
        WHERE id = $3 AND version = $4
        RETURNING id, email, created_at, updated_at`,
		user.Email, user.PasswordHash, id, user.Version,
	).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.Errorf(domain.CONFLICT_ERROR, "conflict updating user")
		}
		return nil, domain.Errorf(domain.INTERNAL_ERROR, "failed to update user: %v", err)
	}

	return &user, nil
}

// Delete takes a user ID and deletes the user from the database.
// It returns an error if the operation fails.
func (s *userService) Delete(id int) error {
	result, err := s.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return domain.Errorf(domain.INTERNAL_ERROR, "error deleting user: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return domain.Errorf(domain.INTERNAL_ERROR, "error getting rows affected: %v", err)
	}

	if rows == 0 {
		return domain.Errorf(domain.NOTFOUND_ERROR, "user not found")
	}

	return nil
}
