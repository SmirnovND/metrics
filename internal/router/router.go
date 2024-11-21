package router

import (
	"github.com/SmirnovND/metrics/internal/controllers"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/SmirnovND/metrics/internal/services/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func Handler(storage *repo.MemStorage, db *sqlx.DB) http.Handler {
	serviceCollector := server.NewCollectorService(storage)
	metricController := controllers.NewMetricsController(serviceCollector)

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Post("/update/{type}/{name}/{value}", metricController.HandleUpdate)
	r.Post("/update", metricController.HandleUpdateJSON)
	r.Post("/updates", metricController.HandleUpdatesJSON)
	r.Get("/value/{type}/{name}", metricController.HandleValue)
	r.Post("/value/{type}/{name}", metricController.HandleValueJSON)
	r.Post("/value", metricController.HandleValueJSON)
	r.Get("/", metricController.HandleRoot)

	healthcheckController := controllers.NewHealthcheckController(db)
	r.Get("/ping", healthcheckController.HandlePing)

	// Обработчик для неподходящего метода (405 Method Not Allowed)
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Обработчик для несуществующих маршрутов (404 Not Found)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Route not found", http.StatusNotFound)
	})

	return r
}
