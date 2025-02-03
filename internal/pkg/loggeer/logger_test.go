package loggeer_test

import (
	"bytes"
	"github.com/SmirnovND/metrics/internal/pkg/loggeer"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithLogging(t *testing.T) {
	// Мокаем логгер для захвата логов в тестах
	var logBuf bytes.Buffer
	log.Logger = log.Output(&logBuf)

	// Создаем тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Эмулируем длительное выполнение
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World"))
	})

	// Оборачиваем обработчик в наш middleware
	loggedHandler := loggeer.WithLogging(handler)

	// Создаем тестовый HTTP-запрос
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// Выполняем запрос
	loggedHandler.ServeHTTP(w, req)

	// Проверяем код статуса ответа
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Проверяем, что лог содержит ожидаемые данные
	loggedOutput := logBuf.String()
	assert.Contains(t, loggedOutput, "Request information")
	assert.Contains(t, loggedOutput, "uri")
	assert.Contains(t, loggedOutput, "method")
	assert.Contains(t, loggedOutput, "status")
	assert.Contains(t, loggedOutput, "size")
	assert.Contains(t, loggedOutput, "duration")

	// Убедитесь, что размер логирования и статус присутствуют
	assert.Contains(t, loggedOutput, "GET")
	assert.Contains(t, loggedOutput, "200")
}
