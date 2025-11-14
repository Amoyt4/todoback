package repository

import (
	"context"
	"diplom_back/internal/entity"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllCleanings(ctx context.Context, db *pgxpool.Pool) ([]entity.Cleanings, error) {
	var cleanings []entity.Cleanings

	rows, err := db.Query(ctx, "SELECT * FROM cleaning")
	if err != nil {
		slog.Warn("GetAllCleanings: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cleaning entity.Cleanings
		err := rows.Scan(
			&cleaning.ID,
			&cleaning.RoomNum,
			&cleaning.StartTime,
			&cleaning.EndTime,
			&cleaning.Comment)
		if err != nil {
			slog.Warn("GetAllCleanings: ", err)
			return nil, err
		}
		cleanings = append(cleanings, cleaning)
	}
	if err := rows.Err(); err != nil {
		slog.Error("GetAllCleanings rows error: ", err)
		return nil, err
	}
	return cleanings, nil
}

func PostNewCleaning(ctx context.Context, db *pgxpool.Pool, cleaning *entity.NewClean) (int, error) {
	var id int
	err := db.QueryRow(ctx, `
        INSERT INTO cleaning (room_num, start_time, end_time, comment)
        VALUES ($1, $2, $3, $4)
        RETURNING id`,
		cleaning.RoomNum, cleaning.StartTime, cleaning.EndTime, cleaning.Comment,
	).Scan(&id)

	return id, err
}

func DeleteCleaning(ctx context.Context, db *pgxpool.Pool, cleaningID uint) error {
	_, err := db.Exec(ctx, "DELETE FROM cleaning WHERE id = $1", cleaningID)
	return err
}
