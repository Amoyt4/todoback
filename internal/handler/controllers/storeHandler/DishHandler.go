package storeHandler

import (
	"context"
	"diplom_back/internal/entity/storeEntity"
	"diplom_back/internal/repository/storeRepository"
	"encoding/json"
	"fmt"

	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllDishesHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dishes, err := storeRepository.GetAllDishes(ctx, db)
		if err != nil {
			slog.Warn("GetAllDishesHandler", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dishes); err != nil {
			slog.Warn("GetAllDishesHandler", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}
	}
}

func PostNewDishesHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var dish storeEntity.NewDish
		body := json.NewDecoder(r.Body)
		if err := body.Decode(&dish); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("JSON не соответствует entity на бэке"))
			return
		}
		defer r.Body.Close()

		if dish.Name == "" || dish.Price == 0 || dish.ImgUrl == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("все поля должны быть заполнены"))
			return
		}

		id, err := storeRepository.PostNewDish(ctx, db, dish)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}
		slog.Info("Успешно добавлено сущность в dishes")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("%d", id)))
		slog.Info("Блюдо успешно добавлено")
	}
}

func DeleteDishByIdHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Content-Type must be application/json"))
			return
		}

		var request storeEntity.DeleteDish

		body := json.NewDecoder(r.Body)
		if err := body.Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		//проверка что пришёл не пустой Id
		if request.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Пришёл пустой id"))
			return
		}
		err := storeRepository.DeleteDish(ctx, db, request.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Dish deleted successfully"))
	}
}
