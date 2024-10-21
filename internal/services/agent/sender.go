package agent

import (
	"bytes"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
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

		// Отправка запроса
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}
		defer resp.Body.Close()

		// Обработка ответа
		if resp.StatusCode == http.StatusOK {
			fmt.Println("Metric sent successfully:", metric.GetName())
		} else {
			fmt.Printf("Failed to send metric %s: %s\n", metric.GetName(), resp.Status)
		}
	}
}
