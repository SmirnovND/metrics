package router

import (
	"github.com/SmirnovND/metrics/internal/controllers"
	"net/http"
)

func Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, controllers.NewMetricsController().HandlePost)
	return mux
}
