package app

import (
	"strings"
	"time"
)

type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateNote struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateNote struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type NoteService interface {
	Create(userID int, req *CreateNote) (*Note, error)
	GetAll(userID int) ([]Note, error)
	Update(userID int, noteID int, req *UpdateNote) (*Note, error)
	Delete(userID int, noteID int) error
}

func (nr *CreateNote) Validate() []string {
	fields := []string{}

	title := strings.TrimSpace(nr.Title)
	if title == "" {
		fields = append(fields, "title is required")
	}
	if len(title) > 200 {
		fields = append(fields, "title must not exceed 200 characters")
	}

	if len(nr.Content) > 10000 {
		fields = append(fields, "content must not exceed 10000 characters")
	}

	if len(fields) == 0 {
		return nil
	}

	return fields
}

func (nr *UpdateNote) Validate() []string {
	fields := []string{}

	if nr.Title == nil && nr.Content == nil {
		fields = append(fields, "update requires at least one field")
		return fields
	}

	if nr.Title != nil {
		title := strings.TrimSpace(*nr.Title)
		if title == "" {
			fields = append(fields, "title is required")
		}
		if len(title) > 200 {
			fields = append(fields, "title must not exceed 200 characters")
		}
	}

	if nr.Content != nil && len(*nr.Content) > 10000 {
		fields = append(fields, "content must not exceed 10000 characters")
	}

	if len(fields) == 0 {
		return nil
	}

	return fields
}
