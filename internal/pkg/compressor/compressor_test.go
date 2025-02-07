package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тестовый обработчик, который просто пишет "OK" в ответ
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func TestWithDecompression(t *testing.T) {
	// Создаём тестовый сервер с middleware
	handler := WithDecompression(http.HandlerFunc(testHandler))
	server := httptest.NewServer(handler)
	defer server.Close()

	// Сжимаем тестовые данные в gzip
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte("test data"))
	if err != nil {
		t.Fatalf("Ошибка при сжатии данных: %v", err)
	}
	gz.Close()

	// Создаём HTTP-запрос с заголовком Content-Encoding: gzip
	req, err := http.NewRequest("POST", server.URL, &buf)
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Encoding", "gzip")

	// Выполняем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, что статус 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	// Проверяем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка при чтении тела ответа: %v", err)
	}
	expected := "OK"
	if string(body) != expected {
		t.Errorf("Ожидалось тело ответа %q, получено %q", expected, string(body))
	}
}

func TestWithCompression(t *testing.T) {
	// Создаём тестовый сервер с middleware
	handler := WithCompression(http.HandlerFunc(testHandler))
	server := httptest.NewServer(handler)
	defer server.Close()

	// Создаём HTTP-запрос с заголовком Accept-Encoding: gzip
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	// Выполняем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, что статус 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}

	// Проверяем, что ответ действительно сжат (заголовок `Content-Encoding: gzip`)
	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("Ожидался заголовок Content-Encoding: gzip, но его нет")
	}

	// Распаковываем тело ответа
	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка при создании gzip-ридера: %v", err)
	}
	defer gzReader.Close()

	body, err := io.ReadAll(gzReader)
	if err != nil {
		t.Fatalf("Ошибка при чтении тела ответа: %v", err)
	}

	// Проверяем, что ответ распаковывается в "OK"
	expected := "OK"
	if string(body) != expected {
		t.Errorf("Ожидалось тело ответа %q, получено %q", expected, string(body))
	}
}
