package controllers

import (
	"github.com/jmoiron/sqlx"
	"net/http"
)

type HealthcheckController struct {
	Db *sqlx.DB
}

func NewHealthcheckController(Db *sqlx.DB) *HealthcheckController {
	return &HealthcheckController{
		Db: Db,
	}
}

func (hc *HealthcheckController) HandlePing(w http.ResponseWriter, r *http.Request) {
	err := hc.Db.Ping()
	if err != nil {
		http.Error(w, "Failed to connect DB", http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
