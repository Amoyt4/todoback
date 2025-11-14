package handler

import (
	"context"
	"diplom_back/config"
	"diplom_back/internal/handler/controllers"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

func Setup(cfg *config.Config, ctx context.Context) http.Handler {
	mux := mux.NewRouter()
	db := cfg.Client

	//Cleaning
	mux.HandleFunc(GetAllCleaning, loggingMiddleware(corsMiddleware(controllers.GetAllCleaningHandler(ctx, db)))).Methods("GET")
	mux.HandleFunc(PostNewCleaning, loggingMiddleware(corsMiddleware(controllers.PostNewCleaningHandler(ctx, db)))).Methods("POST")
	mux.HandleFunc(DeleteCleaningById, loggingMiddleware(corsMiddleware(controllers.DeleteCleaningByIdHandler(ctx, db)))).Methods("DELETE")
	//
	////Employee
	//mux.HandleFunc(GetAllEmployee, loggingMiddleware(corsMiddleware(controllers.GetAllEmployeeHandler(ctx, db)))).Methods("GET")
	//mux.HandleFunc(PostNewEmployee, loggingMiddleware(corsMiddleware(controllers.PostNewEmployeeHandler(ctx, db)))).Methods("POST")
	//mux.HandleFunc(DeleteEmployeeById, loggingMiddleware(corsMiddleware(controllers.DeleteEmployeeByIdHandler(ctx, db)))).Methods("DELETE")
	//
	////Contact Employee
	//mux.HandleFunc(GetAllEmployeeContacts, loggingMiddleware(corsMiddleware(controllers.GetAllEmployeeContactsHandler(ctx, db)))).Methods("GET")
	//mux.HandleFunc(PostNewEmployeeContacts, loggingMiddleware(corsMiddleware(controllers.PostNewEmployeeContactsHandler(ctx, db)))).Methods("POST")
	//mux.HandleFunc(DeleteEmployeeContactsById, loggingMiddleware(corsMiddleware(controllers.DeleteEmployeeContactsByIdHandler(ctx, db)))).Methods("DELETE")
	return mux
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

//func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		allowedOrigins := map[string]bool{
//			"http://localhost:5173":  true,
//			"http://localhost:63342": true,
//		}
//		origin := r.Header.Get("Origin")
//		if allowedOrigins[origin] {
//			//w.Header().Set("Access-Control-Allow-Origin", origin)
//			w.Header().Set("Access-Control-Allow-Credentials", "true")
//		}
//
//		w.Header().Set("Access-Control-Allow-Origin", origin)
//		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
//		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
//
//		if r.Method == "OPTIONS" {
//			w.WriteHeader(http.StatusOK)
//			return
//		}
//
//		next.ServeHTTP(w, r)
//	}
//}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем все origins
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With, X-Auth-Token, Cache-Control, Pragma")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 часа

		// Предварительный запрос (preflight)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
