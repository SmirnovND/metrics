package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
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

		err = baseSend(req, metric, true)
		if err != nil {
			continue
		}
	}
}

// SendJSON Метод для отправки метрик
func SendJSON(m *domain.Metrics, serverHost string) {
	m.Mu.RLock()         // Блокируем доступ к мапе
	defer m.Mu.RUnlock() // Освобождаем доступ после обновления

	for _, metric := range m.Data {
		url := fmt.Sprintf("%s/update/", serverHost)

		// Сериализация метрики в JSON
		jsonData, err := json.Marshal(metric)
		if err != nil {
			fmt.Println("Ошибка при сериализации метрики:", err)
			continue
		}

		// Создание HTTP-запроса
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		err = baseSend(req, metric, true)
		if err != nil {
			continue
		}
	}
}

func baseSend(req *http.Request, metric domain.MetricInterface, enableCompression bool) error {
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
		if metric.GetType() == domain.MetricTypeCounter {
			val := metric.GetValue().(*int64)
			fmt.Println("MetricInterface sent successfully:", metric.GetName(), "Value:", *val)
		}

	} else {
		fmt.Printf("Failed to send metric %s: %s\n", metric.GetName(), resp.Status)
	}

	return nil
}
