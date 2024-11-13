package controllers

import (
	"github.com/SmirnovND/metrics/internal/pkg/paramsparser"
	"github.com/SmirnovND/metrics/internal/services/server"
	serverSaver "github.com/SmirnovND/metrics/internal/usecase/server"
	"net/http"
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
	parseMetric, err := paramsparser.QueryParseMetricAndValue(w, r)
	if err != nil {
		return
	}

	mc.ServiceCollector.SaveMetric(parseMetric)
	w.WriteHeader(http.StatusOK)
}

func (mc *MetricsController) HandleUpdateJSON(w http.ResponseWriter, r *http.Request) {
	parseMetric, err := paramsparser.JSONParseMetric(w, r)
	if err != nil {
		return
	}

	JSONResponse, err := serverSaver.SaveAndFind(parseMetric, mc.ServiceCollector, w)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(JSONResponse)
}

func (mc *MetricsController) HandleUpdatesJSON(w http.ResponseWriter, r *http.Request) {
	parseMetrics, err := paramsparser.JSONParseMetrics(w, r)
	if err != nil {
		return
	}

	JSONResponse, err := serverSaver.SaveAndFindArr(parseMetrics, mc.ServiceCollector, w)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(JSONResponse)
}

func (mc *MetricsController) HandleValue(w http.ResponseWriter, r *http.Request) {
	parseMetric, err := paramsparser.QueryParseMetric(w, r)
	if err != nil {
		return
	}

	metricValue, err := mc.ServiceCollector.GetMetricValue(parseMetric.GetName(), parseMetric.GetType())
	if err != nil {
		http.Error(w, "Not found metric", http.StatusNotFound)
		return
	}

	w.Write([]byte(metricValue))
}

func (mc *MetricsController) HandleValueQueryParamsJSON(w http.ResponseWriter, r *http.Request) {
	parseMetric, err := paramsparser.QueryParseMetric(w, r)
	if err != nil {
		return
	}

	JSONResponse, err := serverSaver.FindAndResponseAsJSON(parseMetric, mc.ServiceCollector, w)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(JSONResponse)
}

func (mc *MetricsController) HandleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", r.Header.Get("Accept"))
	w.WriteHeader(http.StatusOK)
}

func (mc *MetricsController) HandleValueJSON(w http.ResponseWriter, r *http.Request) {
	parseMetric, err := paramsparser.JSONParseMetric(w, r)
	if err != nil {
		return
	}

	JSONResponse, err := serverSaver.FindAndResponseAsJSON(parseMetric, mc.ServiceCollector, w)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(JSONResponse)
}
