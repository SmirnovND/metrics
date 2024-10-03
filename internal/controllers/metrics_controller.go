package controllers

import (
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/services/collector"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type MetricsController struct {
	ServiceCollector *collector.ServiceCollector
}

func NewMetricsController(serviceCollector *collector.ServiceCollector) *MetricsController {
	return &MetricsController{
		ServiceCollector: serviceCollector,
	}
}

func (mc *MetricsController) HandleUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	// Получение параметров из URL
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if metricType == "" {
		http.Error(w, "Invalid URL format", http.StatusNotFound)
		return
	}

	var metric domain.Metric
	if metricType == "gauge" {
		floatValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid Value format", http.StatusBadRequest)
			return
		}
		metric = &domain.Gauge{
			Value: floatValue,
			Name:  metricName,
		}
	} else if metricType == "counter" {
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Value format", http.StatusBadRequest)
			return
		}
		metric = &domain.Counter{
			Value: intValue,
			Name:  metricName,
		}
	} else {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	mc.ServiceCollector.SaveMetric(metric)

	w.WriteHeader(http.StatusOK)
}
func (mc *MetricsController) HandleValue(w http.ResponseWriter, r *http.Request) {
	// Получение параметров из URL
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metric, err := mc.ServiceCollector.FindMetric(metricName, metricType)
	if err != nil {
		http.Error(w, "Not found metric", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	// Определяем тип значения метрики
	switch value := metric.GetValue().(type) {
	case int64:
		// Преобразуем int64 в строку
		w.Write([]byte(fmt.Sprintf("%d", value)))
	case float64:
		// Преобразуем float64 в строку
		w.Write([]byte(fmt.Sprintf("%f", value)))
	default:
		// Если тип не поддерживается, возвращаем ошибку
		http.Error(w, "Unsupported metric type", http.StatusInternalServerError)
	}
}
