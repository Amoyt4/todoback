package storage

import (
	"context"
	"diplom_back/config"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	env := cfg.Env

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		env.DB_USERNAME,
		env.DB_PASSWORD,
		env.DB_HOST,
		env.DB_PORT,
		env.DB_NAME,
	)

	parseConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("Unable to parse parseConfig", err)
	}

	parseConfig.MaxConns = 25
	parseConfig.MinConns = 5
	parseConfig.MaxConnLifetime = 30 * time.Minute
	parseConfig.MaxConnIdleTime = 5 * time.Minute

	conn, err := pgxpool.NewWithConfig(context.Background(), parseConfig)
	if err != nil {
		slog.Error("ошибка при подключении к БД", slog.Any("error", err))
		panic(err)
	}

	// Просто проверяем, что подключение работает
	if err := checkConnection(conn); err != nil {
		slog.Error("ошибка при проверке подключения к БД", slog.Any("error", err))
		panic(err)
	}

	slog.Info("✅ DB connected successfully!")
	return conn
}

func checkConnection(db *pgxpool.Pool) error {
	// Простой запрос для проверки подключения
	var result int
	err := db.QueryRow(context.Background(), "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("не могу выполнить запрос к БД: %w", err)
	}

	log.Println("✅ Подключение к БД работает!")
	return nil
}
