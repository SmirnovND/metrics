package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/pkg/crypto"
	"io"
	"net/http"
)

// Send Метод для отправки метрик
func Send(m *domain.Metrics, serverHost string) {
	m.Mu.RLock()         // Блокируем доступ к мапе
	defer m.Mu.RUnlock() // Освобождаем доступ после обновления

	for _, metric := range m.Data {
		url := fmt.Sprintf("%s/update/%s/%s/%v", serverHost, metric.GetType(), metric.GetName(), metric.GetValue())

		// Создание HTTP-запроса
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		req.Header.Set("Content-Type", "text/plain")

		err = baseSend(req, true)
		if err != nil {
			continue
		}
	}
}

// SendJSON Метод для отправки метрик
func SendJSON(m *domain.Metrics, serverHost string, key string) {
	m.Mu.RLock()         // Блокируем доступ к мапе
	defer m.Mu.RUnlock() // Освобождаем доступ после обновления
	var metrics []*domain.Metric
	url := fmt.Sprintf("%s/updates/", serverHost)

	for _, v := range m.Data {
		metrics = append(metrics, v.(*domain.Metric))
	}

	// Сериализация метрики в JSON
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		fmt.Println("Ошибка при сериализации метрики:", err)
		return
	}

	// Создание HTTP-запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	crypto.SetSignature(req, jsonData, key)
	req.Header.Set("Content-Type", "application/json")

	err = baseSend(req, true)
	if err != nil {
		return
	}
}

func baseSend(req *http.Request, enableCompression bool) error {
	if enableCompression {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)

		_, err := io.Copy(gz, req.Body)
		if err != nil {
			fmt.Println("Ошибка при сжатии запроса:", err)
			return err
		}
		gz.Close()

		// Заменяем тело запроса на сжатое
		req.Body = io.NopCloser(&buf)
		req.ContentLength = int64(buf.Len())
		req.Header.Set("Content-Encoding", "gzip")
	}

	// Отправка запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return err
	}
	defer resp.Body.Close()

	// Обработка ответа
	if resp.StatusCode == http.StatusOK {
		fmt.Println("MetricInterface sent successfully")
	} else {
		fmt.Printf("Failed to send metric")
	}

	return nil
}
