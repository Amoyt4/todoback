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

func GetAllEmployeeContactsHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		contacts, err := repository.GetAllEmployeeContacts(ctx, db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(contacts); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}
	}
}

func PostNewEmployeeContactsHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		var contact entity.NewEmployeeContact

		body := json.NewDecoder(r.Body)
		if err := body.Decode(&contact); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}
		defer r.Body.Close()

		if contact.RoomNum == 0 || contact.Title == "" || contact.Comment == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("%d --- %s --- %s", contact.RoomNum, contact.Title, contact.Comment)))
			return
		}

		id, err := repository.PostNewEmployeeContact(ctx, db, &contact)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		slog.Info("Успешно добавлена сущность EmployeeContact")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("%d - успешно добавоен employee contact", id)))
	}
}

func DeleteEmployeeContactsByIdHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		var request entity.DeleteEmployeeContact
		body := json.NewDecoder(r.Body)
		if err := body.Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}
		defer r.Body.Close()

		if request.Id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Id должен быть заполнен"))
			return
		}

		err := repository.DeleteEmployeeContact(ctx, db, request.Id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%d - удалён employee contact", request.Id)))
	}
}
