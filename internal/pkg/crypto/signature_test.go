package crypto

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockConfig struct {
	key string
}

func (m *MockConfig) GetKey() string {
	return m.key
}

func (m *MockConfig) GetStoreInterval() time.Duration {
	return 10 * time.Second // Примерный фиктивный ответ
}

func (m *MockConfig) GetFileStoragePath() string {
	return "/mock/storage/path" // Примерный фиктивный путь
}

func (m *MockConfig) IsRestore() bool {
	return false // Примерный фиктивный ответ
}

func (m *MockConfig) GetDBDsn() string {
	return "mock_dsn" // Примерный фиктивный DSN
}

func (m *MockConfig) GetFlagRunAddr() string {
	return "http://localhost:8080" // Примерный фиктивный адрес
}

func TestWithHashMiddleware(t *testing.T) {
	// Создание фиктивной конфигурации с ключом для хэширования
	mockConfig := &MockConfig{
		key: "test_key",
	}

	// Создание простого обработчика, который просто отправляет ответ
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
	})

	// Инициализация миддлвара с хэшированием
	handler := WithHashMiddleware(mockConfig, nextHandler)

	// Запуск тестового запроса
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Проверка кода ответа
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}

	// Проверка, что хэш был добавлен в заголовки
	hashHeader := rr.Header().Get("HashSHA256")
	if hashHeader == "" {
		t.Error("Expected HashSHA256 header to be set, but got empty")
	}

	// Проверка, что хэш правильный
	expectedHash := calculateHash([]byte("Hello, world"), mockConfig.GetKey())
	if hashHeader != expectedHash {
		t.Errorf("Expected HashSHA256 header to be %v, but got %v", expectedHash, hashHeader)
	}
}

func TestWithCryptoKeyMiddlewareValidHash(t *testing.T) {
	// Создание фиктивной конфигурации с ключом
	mockConfig := &MockConfig{
		key: "test_key",
	}

	// Тело запроса с данными
	body := []byte("Hello, world")

	// Создание запроса с хэшем в заголовке
	hash := calculateHash(body, mockConfig.GetKey())
	req := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewReader(body))
	req.Header.Set("HashSHA256", hash)

	// Создание простого обработчика
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Received"))
	})

	// Инициализация миддлвара
	handler := WithCryptoKey(mockConfig, nextHandler)

	// Запуск запроса
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Проверка кода ответа
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}

	// Проверка ответа
	if rr.Body.String() != "Received" {
		t.Errorf("Expected body 'Received', but got %v", rr.Body.String())
	}
}

func TestWithCryptoKeyMiddlewareInvalidHash(t *testing.T) {
	// Создание фиктивной конфигурации с ключом
	mockConfig := &MockConfig{
		key: "test_key",
	}

	// Тело запроса с данными
	body := []byte("Hello, world")

	// Создание запроса с некорректным хэшем в заголовке
	invalidHash := "invalidhash"
	req := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewReader(body))
	req.Header.Set("HashSHA256", invalidHash)

	// Создание простого обработчика
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Received"))
	})

	// Инициализация миддлвара
	handler := WithCryptoKey(mockConfig, nextHandler)

	// Запуск запроса
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Проверка кода ответа
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %v, but got %v", http.StatusBadRequest, rr.Code)
	}

	// Проверка, что ответ содержит ошибку
	if rr.Body.String() != "Invalid HashSHA256\n" {
		t.Errorf("Expected body 'Invalid HashSHA256', but got %v", rr.Body.String())
	}
}

func TestSetSignature(t *testing.T) {
	// Создание фиктивной конфигурации с ключом
	mockConfig := &MockConfig{
		key: "test_key",
	}

	// Данные для подписи
	data := []byte("Hello, world")

	// Создание нового запроса
	req := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewReader(data))

	// Устанавливаем подпись
	SetSignature(req, data, mockConfig.GetKey())

	// Проверка, что хэш установлен в заголовке
	hash := req.Header.Get("HashSHA256")
	if hash == "" {
		t.Error("Expected HashSHA256 header to be set, but got empty")
	}

	// Проверка, что хэш правильный
	expectedHash := calculateHash(data, mockConfig.GetKey())
	if hash != expectedHash {
		t.Errorf("Expected HashSHA256 header to be %v, but got %v", expectedHash, hash)
	}
}
