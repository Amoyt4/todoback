package storeRepository

import (
	"context"
	"diplom_back/internal/entity/storeEntity"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllDishes(ctx context.Context, db *pgxpool.Pool) ([]storeEntity.Dish, error) {
	var dishes []storeEntity.Dish

	rows, err := db.Query(ctx, "SELECT * FROM dishes")
	if err != nil {
		slog.Warn("GetAllDishes", err)
		return dishes, err
	}
	defer rows.Close()
	for rows.Next() {
		var dish storeEntity.Dish
		err := rows.Scan(
			&dish.ID,
			&dish.Name,
			&dish.Price,
			&dish.ImgUrl)
		if err != nil {
			slog.Warn("GetAllDishes", err)
			return dishes, err
		}
		dishes = append(dishes, dish)
	}
	if err := rows.Err(); err != nil {
		slog.Warn("GetAllDishes", err)
		return dishes, err
	}
	return dishes, nil
}

func PostNewDish(ctx context.Context, db *pgxpool.Pool, dish storeEntity.NewDish) (uint, error) {
	var id uint
	err := db.QueryRow(ctx, `INSERT INTO dishes (name, price, img_url) 
							VALUES ($1, $2, $3)
 							RETURNING id`, dish.Name, dish.Price, dish.ImgUrl).Scan(&id)

	return id, err
}

func DeleteDish(ctx context.Context, db *pgxpool.Pool, dishId uint) error {
	_, err := db.Exec(ctx, "DELETE FROM dishes WHERE id = $1", dishId)
	return err
}
