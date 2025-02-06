package controllers

import (
	_ "github.com/SmirnovND/metrics/internal/domain"
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

// @HandleUpdate Обновление метрики
// @Description Обновляет метрику по данным формы
// @Tags Update
// @Accept application/x-www-form-urlencoded
// @Produce text/plain
// @Param type formData string true "Тип метрики"
// @Param name formData string true "Название метрики"
// @Param value formData string true "Значение метрики"
// @Success 200
// @Failure 400
// @Router /update/{type}/{name}/{value} [post]
func (mc *MetricsController) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	parseMetric, err := paramsparser.QueryParseMetricAndValue(w, r)
	if err != nil {
		return
	}

	mc.ServiceCollector.SaveMetric(parseMetric)
	w.WriteHeader(http.StatusOK)
}

// @HandleUpdateJson Обновление метрики
// @Description Обновляет метрику на основе JSON тела запроса
// @Tags Update
// @Accept application/json
// @Produce application/json
// @Param body body domain.Metric true "Данные метрики"
// @Success 200 {object} domain.Metric "Обновленная метрика"
// @Failure 400 {string} string "Ошибка в запросе"
// @Router /update [post]
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

// @HandleUpdatesJSON Обновление метрик
// @Description Обновляет массив метрик на основе JSON тела запроса
// @Tags Update
// @Accept application/json
// @Produce application/json
// @Param body body []domain.Metric true "Массив данных метрик"
// @Success 200 {array} []domain.Metric "Массив обновленных метрик"
// @Failure 400 {string} string "Ошибка в запросе"
// @Router /updates [post]
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

// @HandleGetValue Получение значения метрики по типу и имени
// @Description Получает значение метрики на основе типа и имени
// @Tags Update
// @Accept application/json
// @Produce text/plain
// @Param type path string true "Тип метрики"
// @Param name path string true "Название метрики"
// @Success 200 {string} string "Значение метрики"
// @Failure 400 {string} string "Ошибка в запросе"
// @Router /value/{type}/{name} [get]
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

// @HandleValueJSON Получение метрики по данным в теле запроса
// @Description Принимает метрику без значения и возвращает тот же объект метрики с значением
// @Tags Update
// @Accept application/json
// @Produce application/json
// @Param body body domain.Metric true "Данные метрики"
// @Success 200 {object} domain.Metric "Метрика с значением"
// @Failure 400 {string} string "Ошибка в запросе"
// @Router /value [post]
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
