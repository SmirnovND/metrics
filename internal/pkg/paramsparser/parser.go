package paramsparser

import (
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func JSONParseMetric(w http.ResponseWriter, r *http.Request) (domain.MetricInterface, error) {
	var metric *domain.Metric
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&metric)
	if err != nil {
		http.Error(w, "Error decode:"+err.Error(), http.StatusBadRequest)
		return nil, fmt.Errorf("error decode")
	}
	return metric, nil
}

func JSONParseMetrics(w http.ResponseWriter, r *http.Request) ([]*domain.Metric, error) {
	var metrics []*domain.Metric
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&metrics)
	if err != nil {
		http.Error(w, "Error decode:"+err.Error(), http.StatusBadRequest)
		return nil, fmt.Errorf("error decode")
	}

	return metrics, nil
}

func QueryParseMetric(w http.ResponseWriter, r *http.Request) (domain.MetricInterface, error) {
	// Получение параметров из URL
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	var metric domain.MetricInterface
	if metricType == "" {
		http.Error(w, "Invalid URL format", http.StatusNotFound)
		return nil, fmt.Errorf("invalid URL format")
	}

	if metricType == domain.MetricTypeGauge {
		floatValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid Value format", http.StatusBadRequest)
			return nil, fmt.Errorf("invalid Value format")
		}
		metric = (&domain.Gauge{}).SetType(domain.MetricTypeGauge).SetName(metricName).SetValue(&floatValue)
	} else if metricType == domain.MetricTypeCounter {
		intValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Value format", http.StatusBadRequest)
			return nil, fmt.Errorf("invalid Value format")
		}
		metric = (&domain.Counter{}).SetType(domain.MetricTypeCounter).SetName(metricName).SetValue(&intValue)
	} else {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return nil, fmt.Errorf("invalid URL format")
	}

	return metric, nil
}
