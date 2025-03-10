package domain

import (
	"time"

	"example.com/rest/internal/validator"
)

type User struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Version      int       `json:"-" db:"version"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPatch struct {
	Email       *string `json:"email"`
	NewPassword *string `json:"new_password"`
	Password    string  `json:"password"`
}

// validation

func (u *UserCredentials) Validate() error {
	v := validator.New()

	v.NotBlank(u.Email, "email", "email is required")
	v.Email(u.Email, "email", "email is not valid")

	v.NotBlank(u.Password, "password", "password is required")
	v.BetweenRunes(u.Password, 8, 50, "password", "password must be between 8 and 50 characters long")

	return v.Valid()
}

func (u *UserPatch) Validate() error {
	v := validator.New()

	if u.Email == nil && u.NewPassword == nil {
		v.Message("at least one field must be specified")
		return v.Valid()
	}

	if u.Email != nil {
		v.NotBlank(*u.Email, "email", "email is required")
		v.Email(*u.Email, "email", "email is not valid")
	}

	if u.NewPassword != nil {
		v.NotBlank(*u.NewPassword, "new_password", "new password is required")
		v.BetweenRunes(*u.NewPassword, 8, 50, "new_password", "new password must be between 8 and 50 characters long")
	}

	v.NotBlank(u.Password, "password", "password is required")
	v.BetweenRunes(u.Password, 8, 50, "password", "password must be between 8 and 50 characters long")

	return v.Valid()
}
