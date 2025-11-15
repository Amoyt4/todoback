package controllers

import (
	"context"
	"diplom_back/internal/entity"
	"diplom_back/internal/repository"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllEmployeeHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employees, err := repository.GetAllEmployees(ctx, db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(employees); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}

func PostNewEmployeeHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Content-Type must be application/json"))
			return
		}
		var employee entity.NewEmployee

		body := json.NewDecoder(r.Body)
		if err := body.Decode(&employee); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		if employee.Login == "" || employee.Password == "" || employee.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Логин, пароль, имя не могут быть пустыми"))
			return
		}
		id, err := repository.PostNewEmployee(ctx, db, &employee)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		slog.Info("Успешно добавлено сущность в employee")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("%d - успешно добавлен сотрудник", id)))
	}
}

func DeleteEmployeeByIdHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Content-Type must be application/json"))
			return
		}

		var request entity.DeleteEmployee
		body := json.NewDecoder(r.Body)
		if err := body.Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		if request.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("id не может быть пустым"))
			slog.Warn("Пришёл пустой id при удалении employee")
		}

		err := repository.DeleteEmployee(ctx, db, request.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
}
