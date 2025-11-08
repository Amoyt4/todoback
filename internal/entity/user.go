// internal/entity/user.go
package entity

import "time"

type CreateUserRequest struct {
	TelegramID int64  `json:"telegram_id"` // Теперь это id в БД
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name,omitempty"`
}

type UserResponse struct {
	ID        int64     `json:"id"` // Telegram ID (теперь PRIMARY KEY)
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserResponse struct {
	UserID int64  `json:"user_id"`
	Status string `json:"status"`
}
