package handler

import (
	"diplom_back/config"
	"diplom_back/internal/handler/controllers"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func Setup(cfg *config.Config) {
	mux := mux.NewRouter()
	db := cfg.Client

	mux.HandleFunc(GetAllCleaning, loggingMiddleware(corsMiddleware(controllers.GetAllCleaning(db))))
	mux.HandleFunc(GetAllCleaning, loggingMiddleware(corsMiddleware(controllers.PostNewCleaning(db))))
	mux.HandleFunc(GetAllCleaning, loggingMiddleware(corsMiddleware(controllers.DeleteCleaningById(db))))
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Logging every request

		// не логируем метод OPTIONS
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}

		ip := r.Header.Get("X-Forwarded-For")

		userAgent := r.Header.Get("User-Agent")
		slog.Info(fmt.Sprintf("IP: %s, Method: %s, Route: %s, Query: %s, UserAgent: %s, AuthHeader: %s",
			ip, r.Method, r.URL.Path, r.URL.Query(), userAgent, r.Header.Get("Authorization")))

		next(w, r)
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := map[string]bool{
			"http://localhost:5173":  true,
			"http://localhost:63342": true,
		}
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			//w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
