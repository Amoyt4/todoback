package storeRepository

import (
	"context"
	"diplom_back/internal/entity/storeEntity"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllOrders(ctx context.Context, db *pgxpool.Pool) ([]storeEntity.Order, error) {
	var orders []storeEntity.Order

	rows, err := db.Query(ctx, "SELECT id, room_num, time_to_deliver, total_sum FROM orders ORDER BY id DESC")
	if err != nil {
		slog.Warn("GetAllOrders", err)
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var order storeEntity.Order
		err := rows.Scan(
			&order.ID,
			&order.RoomNum,
			&order.TimeToDeliver,
			&order.TotalSum,
		)
		if err != nil {
			slog.Warn("GetAllOrders", err)
			return orders, err
		}

		// Получаем items для каждого заказа
		items, err := GetOrderItemsByOrderID(ctx, db, order.ID)
		if err != nil {
			slog.Warn("GetAllOrders - get items", err)
			return orders, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		slog.Warn("GetAllOrders", err)
		return orders, err
	}
	return orders, nil
}

func GetOrderByID(ctx context.Context, db *pgxpool.Pool, orderID int) (storeEntity.Order, error) {
	var order storeEntity.Order

	err := db.QueryRow(ctx,
		"SELECT id, room_num, time_to_deliver, total_sum FROM orders WHERE id = $1",
		orderID,
	).Scan(&order.ID, &order.RoomNum, &order.TimeToDeliver, &order.TotalSum)
	if err != nil {
		slog.Warn("GetOrderByID", err)
		return order, err
	}

	// Получаем items для заказа
	items, err := GetOrderItemsByOrderID(ctx, db, orderID)
	if err != nil {
		slog.Warn("GetOrderByID - get items", err)
		return order, err
	}
	order.Items = items

	return order, nil
}

func PostNewOrder(ctx context.Context, db *pgxpool.Pool, orderReq storeEntity.CreateOrderRequest) (int, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	// Сначала вычисляем общую сумму
	totalSum, err := calculateOrderTotal(ctx, tx, orderReq.Items)
	if err != nil {
		return 0, err
	}

	// Вставляем заказ
	var orderID int
	err = tx.QueryRow(ctx,
		`INSERT INTO orders (room_num, time_to_deliver, total_sum) 
		 VALUES ($1, $2, $3) 
		 RETURNING id`,
		orderReq.RoomNum, orderReq.TimeToDeliver, totalSum,
	).Scan(&orderID)
	if err != nil {
		slog.Warn("PostNewOrder - insert order", err)
		return 0, err
	}

	// Вставляем элементы заказа
	for _, item := range orderReq.Items {
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items (order_id, dish_id, quantity) 
			 VALUES ($1, $2, $3)`,
			orderID, item.DishID, item.Quantity,
		)
		if err != nil {
			slog.Warn("PostNewOrder - insert item", err)
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Warn("PostNewOrder - commit", err)
		return 0, err
	}

	return orderID, nil
}

func calculateOrderTotal(ctx context.Context, tx pgx.Tx, items []storeEntity.OrderItemCreate) (int, error) {
	total := 0

	for _, item := range items {
		var price int
		err := tx.QueryRow(ctx,
			"SELECT price FROM dishes WHERE id = $1",
			item.DishID,
		).Scan(&price)
		if err != nil {
			slog.Warn("calculateOrderTotal - get dish price", err)
			return 0, err
		}
		total += price * item.Quantity
	}

	return total, nil
}

func DeleteOrder(ctx context.Context, db *pgxpool.Pool, orderID int) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Сначала удаляем элементы заказа
	_, err = tx.Exec(ctx, "DELETE FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		slog.Warn("DeleteOrder - delete items", err)
		return err
	}

	// Затем удаляем сам заказ
	_, err = tx.Exec(ctx, "DELETE FROM orders WHERE id = $1", orderID)
	if err != nil {
		slog.Warn("DeleteOrder - delete order", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Warn("DeleteOrder - commit", err)
		return err
	}

	return nil
}
