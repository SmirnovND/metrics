package controllers

import (
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/SmirnovND/metrics/internal/services/server"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleHandleUpdate() {
	// Инициализация хранилища и контроллера
	collection := make(map[string]domain.MetricInterface)
	storage := repo.NewMetricRepo(collection)
	serviceCollector := server.NewCollectorService(storage)
	controller := &MetricsController{
		ServiceCollector: serviceCollector,
	}

	// Инициализация роутера
	router := chi.NewRouter()
	router.Post("/update/{type}/{name}/{value}", controller.HandleUpdate)

	// Создание тестового запроса
	req := httptest.NewRequest(http.MethodPost, "/update/counter/my_metric/10", nil)
	rr := httptest.NewRecorder()

	// Передача запроса через роутер
	router.ServeHTTP(rr, req)

	fmt.Println(rr.Code) // HTTP статус
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	//
}

func ExampleHandleUpdateJSON() {
	// Инициализация хранилища и контроллера
	collection := make(map[string]domain.MetricInterface)
	storage := repo.NewMetricRepo(collection)
	serviceCollector := server.NewCollectorService(storage)
	controller := &MetricsController{
		ServiceCollector: serviceCollector,
	}

	// Инициализация роутера
	router := chi.NewRouter()
	router.Post("/update", controller.HandleUpdateJSON)

	// Создание JSON тела запроса
	body := `{"id":"my_metric","type":"counter","value":10}`
	req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Передача запроса через роутер
	router.ServeHTTP(rr, req)

	fmt.Println(rr.Code) // HTTP статус
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	// {"id":"my_metric","type":"counter","value":10}
}

func ExampleHandleUpdatesJSON() {
	// Инициализация хранилища и контроллера
	collection := make(map[string]domain.MetricInterface)
	storage := repo.NewMetricRepo(collection)
	serviceCollector := server.NewCollectorService(storage)
	controller := &MetricsController{
		ServiceCollector: serviceCollector,
	}

	// Инициализация роутера
	router := chi.NewRouter()
	router.Post("/updates", controller.HandleUpdatesJSON)

	// Создание JSON тела запроса
	body := `[{"id":"metric1","type":"gauge","value":3.14},{"id":"metric2","type":"counter","value":42}]`
	req := httptest.NewRequest(http.MethodPost, "/updates", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Передача запроса через роутер
	router.ServeHTTP(rr, req)

	fmt.Println(rr.Code) // HTTP статус
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	// [{"id":"metric1","type":"gauge","value":3.14},{"id":"metric2","type":"counter","value":42}]
}

func ExampleHandleValue() {
	// Инициализация хранилища и контроллера
	collection := make(map[string]domain.MetricInterface)
	storage := repo.NewMetricRepo(collection)

	val := int64(42)
	metric := &domain.Metric{}
	metric.SetType("counter").SetValue(val).SetName("my_metric")
	storage.UpdateMetric(metric)

	serviceCollector := server.NewCollectorService(storage)
	controller := &MetricsController{
		ServiceCollector: serviceCollector,
	}

	// Инициализация роутера
	router := chi.NewRouter()
	router.Get("/value/{type}/{name}", controller.HandleValue)

	// Создание тестового запроса
	req := httptest.NewRequest(http.MethodGet, "/value/counter/my_metric", nil)
	rr := httptest.NewRecorder()

	// Передача запроса через роутер
	router.ServeHTTP(rr, req)

	fmt.Println(rr.Code) // HTTP статус
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	// 42
}

func ExampleHandleValueJSON() {
	// Инициализация хранилища и контроллера
	collection := make(map[string]domain.MetricInterface)
	storage := repo.NewMetricRepo(collection)

	val := int64(42)
	metric := &domain.Metric{}
	metric.SetType("counter").SetValue(val).SetName("my_metric")
	storage.UpdateMetric(metric)

	serviceCollector := server.NewCollectorService(storage)
	controller := &MetricsController{
		ServiceCollector: serviceCollector,
	}

	// Инициализация роутера
	router := chi.NewRouter()
	router.Post("/value", controller.HandleValueJSON)

	// Создание JSON тела запроса
	body := `{"id":"my_metric","type":"counter"}`
	req := httptest.NewRequest(http.MethodPost, "/value", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Передача запроса через роутер
	router.ServeHTTP(rr, req)

	fmt.Println(rr.Code) // HTTP статус
	fmt.Println(rr.Body.String())

	// Output:
	// 200
	// {"id":"my_metric","type":"counter","delta":42}
}
