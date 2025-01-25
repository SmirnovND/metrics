package server

import (
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/services/server"
	"net/http"
)

func SaveAndFind(
	parseMetric domain.MetricInterface,
	ServiceCollector *server.ServiceCollector,
	w http.ResponseWriter,
) ([]byte, error) {
	ServiceCollector.SaveMetric(parseMetric)
	return FindAndResponseAsJSON(parseMetric, ServiceCollector, w)
}

func SaveAndFindArr(
	parseMetrics []*domain.Metric,
	ServiceCollector *server.ServiceCollector,
	w http.ResponseWriter,
) ([]byte, error) {
	var metricsResponse []*domain.Metric
	for _, metric := range parseMetrics {
		ServiceCollector.SaveMetric(metric)

		metricResponse, err := ServiceCollector.FindMetric(metric.GetName(), metric.GetType())
		if err != nil {
			http.Error(w, "Not found metric", http.StatusNotFound)
			return nil, fmt.Errorf("not found metric")
		}

		metricsResponse = append(metricsResponse, metricResponse.(*domain.Metric))
	}

	JSONResponse, err := json.Marshal(metricsResponse)
	if err != nil {
		http.Error(w, "Failed to marshal metric to JSON", http.StatusInternalServerError)
		return nil, fmt.Errorf("failed to marshal metric to JSON")
	}
	return JSONResponse, nil
}

func FindAndResponseAsJSON(
	parseMetric domain.MetricInterface,
	ServiceCollector *server.ServiceCollector,
	w http.ResponseWriter,
) ([]byte, error) {
	metricResponse, err := Find(parseMetric, ServiceCollector, w)
	if err != nil {
		return nil, err
	}

	return JSONResponse(metricResponse, w)
}

func Find(
	parseMetric domain.MetricInterface,
	ServiceCollector *server.ServiceCollector,
	w http.ResponseWriter,
) (domain.MetricInterface, error) {
	metricResponse, err := ServiceCollector.FindMetric(parseMetric.GetName(), parseMetric.GetType())
	if err != nil {
		http.Error(w, "Not found metric", http.StatusNotFound)
		return nil, fmt.Errorf("not found metric")
	}

	return metricResponse, nil
}

func JSONResponse(data interface{}, w http.ResponseWriter) ([]byte, error) {
	JSONResponse, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal metric to JSON", http.StatusInternalServerError)
		return nil, fmt.Errorf("failed to marshal metric to JSON")
	}

	return JSONResponse, nil
}
