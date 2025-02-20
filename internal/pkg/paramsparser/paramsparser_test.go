package paramsparser_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/pkg/paramsparser"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONParseMetric(t *testing.T) {
	validMetric := domain.Metric{}
	validMetric.SetType(domain.MetricTypeGauge).SetName("testMetric").SetValue(42.5)

	body, _ := json.Marshal(validMetric)
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	metric, err := paramsparser.JSONParseMetric(w, r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metric.ID != validMetric.ID || metric.MType != validMetric.MType || metric.Value != validMetric.Value {
		t.Errorf("parsed metric does not match expected")
	}
}

func TestJSONParseMetricsInvalidJSON(t *testing.T) {
	body := []byte(`{"invalid":"json"}`)
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()

	_, err := paramsparser.JSONParseMetrics(w, r)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestQueryParseMetricAndValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Создаем RouteContext и добавляем параметры
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("type", domain.MetricTypeGauge)
	routeCtx.URLParams.Add("name", "cpu_usage")
	routeCtx.URLParams.Add("value", "55.5")

	// Встраиваем RouteContext в общий контекст
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	metric, err := paramsparser.QueryParseMetricAndValue(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metric.GetType() != domain.MetricTypeGauge || metric.GetName() != "cpu_usage" {
		t.Errorf("parsed metric does not match expected")
	}
}

func TestQueryParseMetricAndValueInvalidFloat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Создаем RouteContext и добавляем параметры
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("type", domain.MetricTypeGauge)
	routeCtx.URLParams.Add("name", "cpu_usage")
	routeCtx.URLParams.Add("value", "not_a_number")

	// Встраиваем RouteContext в общий контекст запроса
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	_, err := paramsparser.QueryParseMetricAndValue(w, req)
	if err == nil {
		t.Errorf("expected error for invalid float, got nil")
	}
}
