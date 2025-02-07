package paramsparser

import (
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

const (
	maxBitSize = 64
	radix      = 10
)

// JSONParseMetric парсит JSON-тело запроса и преобразует его в структуру Metric.
func JSONParseMetric(w http.ResponseWriter, r *http.Request) (*domain.Metric, error) {
	var metric *domain.Metric
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&metric)
	if err != nil {
		http.Error(w, "Error decode:"+err.Error(), http.StatusBadRequest)
		return nil, fmt.Errorf("error decode")
	}
	return metric, nil
}

// JSONParseMetrics парсит JSON-тело запроса и преобразует его в массив структур Metric.
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

// QueryParseMetricAndValue парсит параметры запроса из URL и создает объект метрики с заданным значением.
func QueryParseMetricAndValue(w http.ResponseWriter, r *http.Request) (domain.MetricInterface, error) {
	var metric domain.MetricInterface

	// Получение параметров из URL
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")
	metricValue := chi.URLParam(r, "value")

	if metricType == "" {
		http.Error(w, "Invalid URL format", http.StatusNotFound)
		return nil, fmt.Errorf("invalid URL format")
	}

	if metricType == domain.MetricTypeGauge {
		floatValue, err := strconv.ParseFloat(metricValue, maxBitSize)
		if err != nil {
			http.Error(w, "Invalid Value format", http.StatusBadRequest)
			return nil, fmt.Errorf("invalid Value format")
		}
		metric = (&domain.Gauge{}).SetType(domain.MetricTypeGauge).SetName(metricName).SetValue(&floatValue)
	}

	if metricType == domain.MetricTypeCounter {
		intValue, err := strconv.ParseInt(metricValue, radix, maxBitSize)
		if err != nil {
			http.Error(w, "Invalid Value format", http.StatusBadRequest)
			return nil, fmt.Errorf("invalid Value format")
		}
		metric = (&domain.Counter{}).SetType(domain.MetricTypeCounter).SetName(metricName).SetValue(&intValue)
	}

	if metric == nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return nil, fmt.Errorf("invalid URL format")
	}

	return metric, nil
}

// QueryParseMetric парсит параметры запроса из URL и создает объект метрики без значения.
func QueryParseMetric(w http.ResponseWriter, r *http.Request) (domain.MetricInterface, error) {
	var metric domain.MetricInterface

	// Получение параметров из URL
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	if metricType == "" {
		http.Error(w, "Invalid URL format", http.StatusNotFound)
		return nil, fmt.Errorf("invalid URL format")
	}

	if metricType == domain.MetricTypeGauge {
		metric = (&domain.Gauge{}).SetType(domain.MetricTypeGauge).SetName(metricName)
	}

	if metricType == domain.MetricTypeCounter {
		metric = (&domain.Counter{}).SetType(domain.MetricTypeCounter).SetName(metricName)
	}

	if metric == nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return nil, fmt.Errorf("invalid URL format")
	}

	return metric, nil
}
