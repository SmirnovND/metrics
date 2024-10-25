package router

import (
	"github.com/SmirnovND/metrics/internal/controllers"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/SmirnovND/metrics/internal/services/server"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Handler(storage *repo.MemStorage) http.Handler {
	serviceCollector := server.NewCollectorService(storage)
	metricController := controllers.NewMetricsController(serviceCollector)

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", metricController.HandleUpdate)
	r.Post("/update/", metricController.HandleUpdateJson)
	r.Get("/value/{type}/{name}", metricController.HandleValue)
	r.Get("/value/", metricController.HandleValue)

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
