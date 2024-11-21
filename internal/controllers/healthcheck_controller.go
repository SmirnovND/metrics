package controllers

import (
	"github.com/jmoiron/sqlx"
	"net/http"
)

type HealthcheckController struct {
	DB *sqlx.DB
}

func NewHealthcheckController(DB *sqlx.DB) *HealthcheckController {
	return &HealthcheckController{
		DB: DB,
	}
}

func (hc *HealthcheckController) HandlePing(w http.ResponseWriter, r *http.Request) {
	err := hc.DB.Ping()
	if err != nil {
		http.Error(w, "Failed to connect DB", http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
