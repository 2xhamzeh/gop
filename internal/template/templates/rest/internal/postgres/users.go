package postgres

import (
	"context"
	"database/sql"

	"example.com/rest/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Insert takes an email and a password hash and inserts a new user into the database.
// It returns the ID of the created user or an error if the operation fails.
// If the user already exists, it returns an ErrConflict.
func (r *UserRepo) Insert(ctx context.Context, email, passwordHash string) (int, error) {
	var id int
	err := r.db.QueryRowxContext(ctx, "INSERT INTO users (email, password_hash) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id", email, passwordHash).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrConflict
		}
		return 0, err
	}
	return id, nil
}

// GetByID takes a user ID and finds the user in the database.
// It returns the user or an error if the operation fails.
// If the user is not found, it returns an ErrNotFound.
func (r *UserRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRowxContext(ctx, "SELECT * FROM users WHERE id = $1", id).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail takes an email and finds the user in the database.
// It returns the user or an error if the operation fails.
// If the user is not found, it returns an ErrNotFound.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRowxContext(ctx, "SELECT * FROM users WHERE email = $1", email).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update takes a user object and updates the user in the database overwriting the existing user.
// It returns the updated user or an error if the operation fails.
// If the user is not found (based on the user ID and version), it returns an ErrNotFound.
func (r *UserRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	// update the user with the new user object
	var updated domain.User
	err := r.db.QueryRowxContext(ctx, "UPDATE users SET email = $1, password_hash = $2, updated_at = now(), version = version + 1 WHERE id = $3 AND version = $4 RETURNING *", user.Email, user.PasswordHash, user.ID, user.Version).StructScan(&updated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &updated, nil
}

// Delete takes a user ID and deletes the user from the database.
// It returns an error if the operation fails.
// If the user is not found, it returns an ErrNotFound.
func (r *UserRepo) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
