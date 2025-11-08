// internal/adapter/postgres/note.go
package postgres

import (
	"context"
	"diplom_back/internal/entity"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateNote создает новую заметку
func CreateNote(ctx context.Context, db *pgxpool.Pool, userID int64, title, content string) (int64, error) {
	var noteID int64
	err := db.QueryRow(ctx,
		"INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id",
		userID, title, content,
	).Scan(&noteID)

	if err != nil {
		return 0, fmt.Errorf("failed to create note: %w", err)
	}

	return noteID, nil
}

// GetUserNotes возвращает все заметки пользователя
func GetUserNotes(ctx context.Context, db *pgxpool.Pool, userID int64) ([]entity.NoteResponse, error) {
	rows, err := db.Query(ctx,
		`SELECT n.id, n.user_id, n.title, n.content, n.created_at, 
		 u.username, u.first_name 
		 FROM notes n 
		 LEFT JOIN users u ON n.user_id = u.id 
		 WHERE n.user_id = $1 
		 ORDER BY n.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notes: %w", err)
	}
	defer rows.Close()

	var notes []entity.NoteResponse
	for rows.Next() {
		var note entity.NoteResponse
		var username, firstName *string

		err := rows.Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &username, &firstName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}

		// Добавляем информацию о пользователе
		if username != nil {
			note.Username = *username
		}
		if firstName != nil {
			note.FirstName = *firstName
		}

		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return notes, nil
}

// DeleteNote удаляет заметку
func DeleteNote(ctx context.Context, db *pgxpool.Pool, noteID int64) error {
	result, err := db.Exec(ctx,
		"DELETE FROM notes WHERE id = $1",
		noteID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	// Проверяем, что заметка была удалена
	if result.RowsAffected() == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}
