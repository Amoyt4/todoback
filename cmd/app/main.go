package main

import (
	"context"
	"diplom_back/config"
	handler "diplom_back/internal/handler/http"
	"diplom_back/internal/storage"
	"errors"
	"log/slog"
	"net/http"
	"os" // 🔥 ДОБАВЬТЕ ЭТОТ ИМПОРТ
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.GetConfig()

	slog.Info("Starting app")
	slog.Debug("Debud messages are enabled")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg.Client = storage.NewConnection(ctx, cfg)

	// 🔥 ИСПРАВЛЕНИЕ: Получаем порт из переменных окружения Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // дефолтный порт для Render
	}

	server := &http.Server{
		Addr:         ":" + port, // 🔥 Используем порт из Render
		Handler:      handler.Setup(cfg, ctx),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		slog.Info("Server running on port " + port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", slog.String("error", err.Error()))
			panic(err)
		}
	}()

	<-ctx.Done()
	slog.Info("Graceful shutdown initiated...")
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
		panic(err)
	}
}
