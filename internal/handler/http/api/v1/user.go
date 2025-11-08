package v1

import (
	"context"
	"diplom_back/internal/adapter/postgres"
	"diplom_back/internal/entity"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUserHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req entity.CreateUserRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error:   "invalid_json",
				Message: "Invalid JSON in request body",
			})
			return
		}

		// TelegramID теперь напрямую становится id в БД
		userID, err := postgres.GetOrCreateUser(r.Context(), db, req.TelegramID, req.Username, req.FirstName, req.LastName)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.ErrorResponse{
				Error:   "create_user_failed",
				Message: "Failed to create user: " + err.Error(),
			})
			return
		}

		response := entity.CreateUserResponse{
			UserID: userID,
			Status: "success",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
