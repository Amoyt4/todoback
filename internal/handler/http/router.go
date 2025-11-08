package handler

import (
	"context"
	"diplom_back/config"
	v1 "diplom_back/internal/handler/http/api/v1"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

func Setup(cfg *config.Config, ctx context.Context) http.Handler {
	mux := http.NewServeMux()
	db := cfg.Client

	// Корневой эндпоинт для проверки здоровья
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Backend is running!",
			"time":    time.Now().String(),
		})
	})

	// 🔥 ИСПРАВЛЕНИЕ: Правильное применение middleware
	// Сначала CORS, потом логирование, потом хендлер
	mux.Handle(CreateUser, corsMiddleware(loggingMiddleware(v1.CreateUserHandler(ctx, db))))
	mux.Handle(CreateNote, corsMiddleware(loggingMiddleware(v1.CreateNoteHandler(ctx, db))))
	mux.Handle(GetNotes, corsMiddleware(loggingMiddleware(v1.GetNotesHandler(ctx, db))))
	mux.Handle(DeleteNote, corsMiddleware(loggingMiddleware(v1.DeleteNoteHandler(ctx, db))))

	return mux
}

// 🔥 ИСПРАВЛЕНИЕ: loggingMiddleware теперь принимает http.Handler
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логируем OPTIONS запросы тоже
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}

		userAgent := r.Header.Get("User-Agent")
		slog.Info("HTTP request",
			"ip", ip,
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"userAgent", userAgent,
		)

		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем все origins (для тестирования)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Обрабатываем preflight OPTIONS запрос
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
