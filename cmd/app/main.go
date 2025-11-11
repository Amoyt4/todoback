package main

import (
	"diplom_back/config"
	"diplom_back/internal/handler"
	"diplom_back/internal/storage"
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	cfg := config.GetConfig()

	slog.Info("Starting application")

	cfg.Client = storage.NewConnection(cfg)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Env.API_PORT),
		Handler: handler.Setup(cfg),
	}

}
