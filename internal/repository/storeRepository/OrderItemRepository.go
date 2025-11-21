package storeRepository

import (
	"context"
	"diplom_back/internal/entity/storeEntity"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetOrderItemsByOrderID(ctx context.Context, db *pgxpool.Pool, orderID int) ([]storeEntity.OrderItem, error) {
	var items []storeEntity.OrderItem

	rows, err := db.Query(ctx,
		"SELECT id, order_id, dish_id, quantity FROM order_items WHERE order_id = $1 ORDER BY id",
		orderID,
	)
	if err != nil {
		slog.Warn("GetOrderItemsByOrderID", err)
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item storeEntity.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.DishID,
			&item.Quantity,
		)
		if err != nil {
			slog.Warn("GetOrderItemsByOrderID - scan", err)
			return items, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		slog.Warn("GetOrderItemsByOrderID - rows error", err)
		return items, err
	}

	return items, nil
}

func UpdateOrderItemQuantity(ctx context.Context, db *pgxpool.Pool, orderItemID int, quantity int) error {
	_, err := db.Exec(ctx,
		"UPDATE order_items SET quantity = $1 WHERE id = $2",
		quantity, orderItemID,
	)
	if err != nil {
		slog.Warn("UpdateOrderItemQuantity", err)
		return err
	}
	return nil
}

func DeleteOrderItem(ctx context.Context, db *pgxpool.Pool, orderItemID int) error {
	_, err := db.Exec(ctx,
		"DELETE FROM order_items WHERE id = $1",
		orderItemID,
	)
	if err != nil {
		slog.Warn("DeleteOrderItem", err)
		return err
	}
	return nil
}
