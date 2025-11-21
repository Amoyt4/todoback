package storeHandler

import (
	"context"
	"diplom_back/internal/entity/storeEntity"
	"diplom_back/internal/repository/storeRepository"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetOrderItemsByOrderIDHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Content-Type must be application/json"))
			return
		}

		var request storeEntity.GetOrderItemsRequest
		body := json.NewDecoder(r.Body)
		if err := body.Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("JSON не соответствует entity на бэке"))
			return
		}
		defer r.Body.Close()

		if request.OrderID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Пришёл пустой order_id"))
			return
		}

		items, err := storeRepository.GetOrderItemsByOrderID(ctx, db, request.OrderID)
		if err != nil {
			slog.Warn("GetOrderItemsByOrderIDHandler", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(items); err != nil {
			slog.Warn("GetOrderItemsByOrderIDHandler", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"` + err.Error() + `"}`))
			return
		}
	}
}

func UpdateOrderItemQuantityHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Content-Type must be application/json"))
			return
		}

		var request storeEntity.UpdateOrderItemRequest
		body := json.NewDecoder(r.Body)
		if err := body.Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		if request.OrderItemID == 0 || request.Quantity == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("order_item_id и quantity обязательны"))
			return
		}

		err := storeRepository.UpdateOrderItemQuantity(ctx, db, request.OrderItemID, request.Quantity)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Order item quantity updated successfully"))
	}
}

func DeleteOrderItemHandler(ctx context.Context, db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Content-Type must be application/json"))
			return
		}

		var request storeEntity.DeleteOrderItemRequest
		body := json.NewDecoder(r.Body)
		if err := body.Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()

		if request.OrderItemID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Пришёл пустой order_item_id"))
			return
		}

		err := storeRepository.DeleteOrderItem(ctx, db, request.OrderItemID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Order item deleted successfully"))
	}
}
