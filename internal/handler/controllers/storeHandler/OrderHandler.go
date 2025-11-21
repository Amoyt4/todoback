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

func GetAllOrdersHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := storeRepository.GetAllOrders(ctx, db)
		if err != nil {
			slog.Warn("GetAllOrdersHandler", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			slog.Warn("GetAllOrdersHandler", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}
	}
}

func GetOrderByIDHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Content-Type must be application/json"))
			return
		}

		var request storeEntity.GetOrderByIDRequest
		body := json.NewDecoder(r.Body)
		if err := body.Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("JSON не соответствует entity на бэке"))
			return
		}
		defer r.Body.Close()

		if request.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Пришёл пустой id"))
			return
		}

		order, err := storeRepository.GetOrderByID(ctx, db, request.ID)
		if err != nil {
			slog.Warn("GetOrderByIDHandler", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			slog.Warn("GetOrderByIDHandler", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}
	}
}

func PostNewOrderHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var orderReq storeEntity.CreateOrderRequest
		body := json.NewDecoder(r.Body)
		if err := body.Decode(&orderReq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("JSON не соответствует entity на бэке"))
			return
		}
		defer r.Body.Close()

		// Валидация
		if orderReq.RoomNum == 0 || orderReq.TimeToDeliver.IsZero() || len(orderReq.Items) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("все поля должны быть заполнены и должен быть хотя бы один item"))
			return
		}

		// Проверяем каждый item
		for _, item := range orderReq.Items {
			if item.DishID == 0 || item.Quantity == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("все items должны иметь dish_id и quantity"))
				return
			}
		}

		id, err := storeRepository.PostNewOrder(ctx, db, orderReq)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}

		slog.Info("Успешно создан заказ")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"id":%d}`, id)))
	}
}

func DeleteOrderHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Content-Type must be application/json"))
			return
		}

		var request storeEntity.DeleteOrder
		body := json.NewDecoder(r.Body)
		if err := body.Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		if request.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Пришёл пустой id"))
			return
		}

		err := storeRepository.DeleteOrder(ctx, db, request.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Order deleted successfully"))
	}
}
