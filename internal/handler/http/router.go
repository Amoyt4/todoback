package handler

import (
	"context"
	"diplom_back/config"
	v1 "diplom_back/internal/handler/http/api/v1"
	"log/slog"
	"net/http"
)

func Setup(cfg *config.Config, ctx context.Context) http.Handler {
	mux := http.NewServeMux()
	db := cfg.Client
	mux.HandleFunc(CreateUser, loggingMiddleware(corsMiddleware(v1.CreateUserHandler(ctx, db))))
	mux.HandleFunc(CreateNote, loggingMiddleware(corsMiddleware(v1.CreateNoteHandler(ctx, db))))
	mux.HandleFunc(GetNotes, loggingMiddleware(corsMiddleware(v1.GetNotesHandler(ctx, db))))
	mux.HandleFunc(DeleteNote, loggingMiddleware(corsMiddleware(v1.DeleteNoteHandler(ctx, db))))

	return mux
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}

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

		next(w, r)
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
