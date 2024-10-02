package controllers

import (
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/usecase"
	"net/http"
	"strconv"
	"strings"
)

type MetricsController struct{}

func NewMetricsController() *MetricsController {
	return &MetricsController{}
}

func (mc *MetricsController) HandleUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 5 || parts[2] == "" {
		http.Error(w, "Invalid URL format", http.StatusNotFound)
		return
	}

	var metric domain.Metric
	if parts[2] == "gauge" {
		floatValue, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			http.Error(w, "Invalid Value format", http.StatusBadRequest)
			return
		}
		metric = &domain.Gauge{
			Value: floatValue,
			Name:  parts[3],
		}
	} else if parts[2] == "counter" {
		intValue, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			http.Error(w, "Invalid Value format", http.StatusBadRequest)
			return
		}
		metric = &domain.Counter{
			Value: intValue,
			Name:  parts[3],
		}
	} else {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	usecase.ProcessMetrics(metric)

	w.WriteHeader(http.StatusOK)
}
