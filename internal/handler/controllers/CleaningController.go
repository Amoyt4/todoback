package controllers

import (
	"diplom_back/config"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllCleaning(db *pgxpool.Pool) http.HandlerFunc {

}

func PostNewCleaning(db *config.Config) http.HandlerFunc {

}

func DeleteCleaningById(db *config.Config) http.HandlerFunc {

}
