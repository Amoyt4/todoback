// internal/handler/http/api/v1/notes.go
package v1

import (
	"context"
	"diplom_back/internal/adapter/postgres"
	"diplom_back/internal/entity"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateNoteHandler - создание новой заметки
func CreateNoteHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req entity.CreateNoteRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error:   "invalid_json",
				Message: "Invalid JSON in request body",
			})
			return
		}

		noteID, err := postgres.CreateNote(r.Context(), db, req.UserID, req.Title, req.Content)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error:   "create_note_failed",
				Message: "Failed to create note: " + err.Error(),
			})
			return
		}

		response := entity.CreateNoteResponse{
			NoteID: noteID,
			Status: "created",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// GetNotesHandler - все заметки пользователя (по query параметру user_id)
func GetNotesHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("user_id")
		if userIDStr == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error:   "missing_user_id",
				Message: "user_id parameter is required",
			})
			return
		}

		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error:   "invalid_user_id",
				Message: "Invalid user_id parameter",
			})
			return
		}

		notes, err := postgres.GetUserNotes(r.Context(), db, userID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error:   "get_notes_failed",
				Message: "Failed to get user notes: " + err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
	}
}

// DeleteNoteHandler - удаление заметки
func DeleteNoteHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем id из пути: /api/v1/notes/123
		path := r.URL.Path
		parts := strings.Split(path, "/")

		if len(parts) < 5 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error:   "invalid_path",
				Message: "Invalid path format",
			})
			return
		}

		noteIDStr := parts[4]
		noteID, err := strconv.ParseInt(noteIDStr, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error:   "invalid_note_id",
				Message: "Invalid note ID parameter",
			})
			return
		}

		err = postgres.DeleteNote(r.Context(), db, noteID)
		if err != nil {
			if err.Error() == "note not found" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(entity.ErrorResponse{
					Error:   "note_not_found",
					Message: "Note with specified ID not found",
				})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error:   "delete_note_failed",
				Message: "Failed to delete note: " + err.Error(),
			})
			return
		}

		response := entity.DeleteNoteResponse{
			Status: "deleted",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
