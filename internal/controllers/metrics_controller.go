package controllers

import (
	"encoding/json"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/services/server"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type MetricsController struct {
	ServiceCollector *server.ServiceCollector
}

func NewMetricsController(serviceCollector *server.ServiceCollector) *MetricsController {
	return &MetricsController{
		ServiceCollector: serviceCollector,
	}
}

func (mc *MetricsController) HandleUpdate(w http.ResponseWriter, r *http.Request) {

	// Получение параметров из URL
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if metricType == "" {
		http.Error(w, "Invalid URL format", http.StatusNotFound)
		return
	}

	var metric domain.MetricInterface
	if metricType == domain.MetricTypeGauge {
		floatValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid Value format", http.StatusBadRequest)
			return
		}
		metric = (&domain.Gauge{}).SetType(domain.MetricTypeGauge).SetName(metricName).SetValue(&floatValue)
	} else if metricType == domain.MetricTypeCounter {
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Value format", http.StatusBadRequest)
			return
		}
		metric = (&domain.Counter{}).SetType(domain.MetricTypeCounter).SetName(metricName).SetValue(&intValue)
	} else {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	mc.ServiceCollector.SaveMetric(metric)

	w.WriteHeader(http.StatusOK)
}

func (mc *MetricsController) HandleUpdateJSON(w http.ResponseWriter, r *http.Request) {

	var metric *domain.Metric
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&metric)
	if err != nil {
		http.Error(w, "Error decode:"+err.Error(), http.StatusBadRequest)
		return
	}

	mc.ServiceCollector.SaveMetric(metric)

	metricResponse, err := mc.ServiceCollector.FindMetric(metric.GetName(), metric.GetType())
	if err != nil {
		http.Error(w, "Not found metric", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(metricResponse)
	if err != nil {
		http.Error(w, "Failed to marshal metric to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (mc *MetricsController) HandleValue(w http.ResponseWriter, r *http.Request) {
	// Получение параметров из URL
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue, err := mc.ServiceCollector.GetMetricValue(metricName, metricType)
	if err != nil {
		http.Error(w, "Not found metric", http.StatusNotFound)
		return
	}

	w.Write([]byte(metricValue))
}

func (mc *MetricsController) HandleValueQueryParamsJSON(w http.ResponseWriter, r *http.Request) {
	// Получение параметров из URL
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	metricResponse, err := mc.ServiceCollector.FindMetric(metricName, metricType)
	if err != nil {
		http.Error(w, "Not found metric", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(metricResponse)
	if err != nil {
		http.Error(w, "Failed to marshal metric to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (mc *MetricsController) HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (mc *MetricsController) HandleValueJSON(w http.ResponseWriter, r *http.Request) {
	var metric *domain.Metric
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&metric)
	if err != nil {
		metric = &domain.Metric{}
		metric.SetName(chi.URLParam(r, "name"))
		metric.SetType(chi.URLParam(r, "type"))
		//http.Error(w, "Error decode:"+err.Error(), http.StatusBadRequest)
		//return
	}

	metricResponse, err := mc.ServiceCollector.FindMetric(metric.GetName(), metric.GetType())
	if err != nil {
		http.Error(w, "Not found metric", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(metricResponse)
	if err != nil {
		http.Error(w, "Failed to marshal metric to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
