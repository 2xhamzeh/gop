package app

import (
	"strings"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUser struct {
	Username *string `json:"username"`
}

type UserService interface {
	Create(req *UserCredentials) (*User, error)
	Get(id int) (*User, error)
	GetByUsername(username string) (*User, error)
	Authenticate(req *UserCredentials) (*User, error)
	Update(id int, req *UpdateUser) (*User, error)
	Delete(id int) error
}

func (u *UserCredentials) Validate() []string {
	fields := []string{}

	username := strings.TrimSpace(u.Username)
	if username == "" {
		fields = append(fields, "username is required")
	}
	if len(username) < 3 || len(username) > 50 {
		fields = append(fields, "username must be between 3 and 50 characters")
	}

	if u.Password == "" {
		fields = append(fields, "password is required")
	}
	if len(u.Password) < 8 || len(u.Password) > 50 {
		fields = append(fields, "password must be between 8 and 50 characters")
	}

	if len(fields) == 0 {
		return nil
	}

	return fields
}

func (u *UpdateUser) Validate() []string {
	fields := []string{}

	if u.Username == nil {
		fields = append(fields, "update requires at least one field")
		return fields
	}

	if u.Username != nil {
		username := strings.TrimSpace(*u.Username)
		if len(username) < 3 || len(username) > 50 {
			fields = append(fields, "username must be between 3 and 50 characters")
		}
	}

	if len(fields) == 0 {
		return nil
	}

	return fields
}
