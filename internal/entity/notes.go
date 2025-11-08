// internal/entity/note.go
package entity

import "time"

type CreateNoteRequest struct {
	UserID  int64  `json:"user_id"` // Telegram ID пользователя
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	// Дополнительная информация о пользователе
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
}

type CreateNoteResponse struct {
	NoteID int64  `json:"note_id"`
	Status string `json:"status"`
}

type DeleteNoteResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
