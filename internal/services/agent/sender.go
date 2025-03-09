package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/pkg/crypto"
	"github.com/SmirnovND/metrics/internal/pkg/system"
	"github.com/SmirnovND/metrics/pb"
	"io"
	"net/http"
	"time"
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

type Sender interface {
	Send(metrics []*domain.Metric) error
}

type HTTPSender struct {
	ServerHost string
	Key        string
}

// SendJSON метод для отправки метрик с шифрованием
func (h *HTTPSender) Send(metrics []*domain.Metric) error {
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации метрики: %w", err)
	}

	url := fmt.Sprintf("%s/updates/", h.ServerHost)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	crypto.SetSignature(req, jsonData, h.Key)
	req.Header.Set("Content-Type", "application/json")

	return baseSend(req, true)
}

func SendJSON(m *domain.Metrics, sender Sender) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()

	var metrics []*domain.Metric
	for _, v := range m.Data {
		metrics = append(metrics, v)
	}

	err := sender.Send(metrics)
	if err != nil {
		fmt.Println("Ошибка при отправке метрик:", err)
	} else {
		fmt.Println("Метрики успешно отправлены")
	}
}

type GRPCSender struct {
	Client pb.MetricsServiceClient
}

func (g *GRPCSender) Send(metrics []*domain.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.MetricsRequest{}
	for _, m := range metrics {
		req.Metrics = append(req.Metrics, &pb.Metric{
			Id:    m.ID,
			Value: m.Value,
		})
	}

	_, err := g.Client.SendMetrics(ctx, req)
	return err
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

		ip, _ := system.GetLocalIP()
		req.Header.Set("X-Real-IP", ip)
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
		return errors.New("Failed to send metric")
	}

	return nil
}
