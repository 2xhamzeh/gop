package domain

import (
	"regexp"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Version      int       `json:"-"`
}

type UserCredentials struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

type UserPatch struct {
	Email       *string `json:"email"`
	NewPassword *string `json:"new_password"`
	Password    string  `json:"password"`
}

// validation

var (
	RgxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func isEmail(value string) bool {
	if len(value) > 254 {
		return false
	}

	return RgxEmail.MatchString(value)
}

func (u *UserCredentials) Validate() []string {
	fields := []string{}

	if u.Email == "" {
		fields = append(fields, "email is required")
	} else if !isEmail(u.Email) {
		fields = append(fields, "invalid email")
	}

	if u.Password == "" {
		fields = append(fields, "password is required")
	} else if len(u.Password) < 8 || len(u.Password) > 50 {
		fields = append(fields, "password must be between 8 and 50 characters")
	}

	if len(fields) == 0 {
		return nil
	}

	return fields
}

func (u *UserPatch) Validate() []string {
	fields := []string{}

	if u.Email == nil && u.NewPassword == nil {
		fields = append(fields, "no fields to patch")
		return fields
	}

	if u.Email != nil {
		if !isEmail(*u.Email) {
			fields = append(fields, "invalid email")
		}
	}

	if u.NewPassword != nil {
		if len(*u.NewPassword) < 8 || len(*u.NewPassword) > 50 {
			fields = append(fields, "new password must be between 8 and 50 characters")
		}
	}

	if u.Password == "" {
		fields = append(fields, "password is required")
	} else if len(u.Password) < 8 || len(u.Password) > 50 {
		fields = append(fields, "password must be between 8 and 50 characters")
	}

	if len(fields) == 0 {
		return nil
	}

	return fields
}
