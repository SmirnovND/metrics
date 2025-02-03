package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SmirnovND/metrics/internal/middleware"
	"github.com/stretchr/testify/assert"
)

// Простые middleware-функции для тестирования
func addHeaderMiddleware(header, value string) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(header, value) // Добавляем заголовок
			next.ServeHTTP(w, r)
		})
	}
}

func statusCodeMiddleware(statusCode int) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(statusCode) // Устанавливаем статус код
			next.ServeHTTP(w, r)
		})
	}
}

func TestChainMiddleware(t *testing.T) {
	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Основной код обработки
		w.Write([]byte("Hello, World"))
	})

	// Создаем цепочку middleware
	chain := middleware.ChainMiddleware(handler,
		addHeaderMiddleware("X-Test-Header", "HeaderValue"),
		statusCodeMiddleware(http.StatusAccepted),
	)

	// Создаем тестовый запрос
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// Выполняем запрос через цепочку middleware
	chain.ServeHTTP(w, req)

	// Проверяем, что статус код ответа правильный
	assert.Equal(t, http.StatusAccepted, w.Result().StatusCode)

	// Проверяем, что заголовок был добавлен
	assert.Equal(t, "HeaderValue", w.Header().Get("X-Test-Header"))

	// Проверяем, что тело ответа правильно записано
	assert.Equal(t, "Hello, World", w.Body.String())
}
