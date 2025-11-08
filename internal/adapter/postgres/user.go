// internal/adapter/postgres/user.go
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetOrCreateUser теперь создает/возвращает пользователя по telegram_id (который теперь id)
func GetOrCreateUser(ctx context.Context, db *pgxpool.Pool, telegramID int64, username, firstName, lastName string) (int64, error) {
	var userID int64

	// Пробуем найти пользователя по telegram_id (который теперь id)
	err := db.QueryRow(ctx,
		"SELECT id FROM users WHERE id = $1",
		telegramID,
	).Scan(&userID)

	if err != nil {
		// Пользователя нет - создаем нового
		err = db.QueryRow(ctx,
			"INSERT INTO users (id, username, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id",
			telegramID, username, firstName, lastName,
		).Scan(&userID)

		if err != nil {
			return 0, fmt.Errorf("failed to create user: %w", err)
		}
	}

	return userID, nil
}

// GetUserIDByTelegramID теперь просто возвращает telegram_id (так как он теперь id)
func GetUserIDByTelegramID(ctx context.Context, db *pgxpool.Pool, telegramID int64) (int64, error) {
	var userID int64
	err := db.QueryRow(ctx,
		"SELECT id FROM users WHERE id = $1",
		telegramID,
	).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}

	return userID, nil
}
