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

// SendJSON метод для отправки метрик с шифрованием
func SendJSON(m *domain.Metrics, serverHost string, key string) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	var metrics []*domain.Metric
	url := fmt.Sprintf("%s/updates/", serverHost)

	for _, v := range m.Data {
		metrics = append(metrics, v)
	}

	// Сериализация метрики в JSON
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		fmt.Println("Ошибка при сериализации метрики:", err)
		return
	}

	// 2. Генерация AES-ключа
	aesKey, err := crypto.GenerateAESKey()
	if err != nil {
		fmt.Println("Ошибка генерации AES-ключа:", err)
		return
	}

	// 3. Шифрование JSON с помощью AES
	encryptedData, err := crypto.EncryptAES(jsonData, aesKey)
	if err != nil {
		fmt.Println("Ошибка шифрования AES:", err)
		return
	}

	// 4. Загрузка публичного RSA-ключа
	publicKey, err := crypto.LoadPublicKey(key)
	if err != nil {
		fmt.Println("Ошибка загрузки публичного ключа:", err)
		return
	}

	// 5. Шифрование AES-ключа с помощью RSA
	encryptedAESKey, err := crypto.EncryptAESKey(aesKey, publicKey)
	if err != nil {
		fmt.Println("Ошибка шифрования AES-ключа:", err)
		return
	}

	// 6. Формирование запроса
	requestPayload := map[string]string{
		"key":  encryptedAESKey,
		"data": encryptedData,
	}

	finalJSON, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Println("Ошибка сериализации запроса:", err)
		return
	}

	// Создание HTTP-запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(finalJSON))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

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
		fmt.Printf("Failed to send metric: code - %d \n", resp.StatusCode)
	}

	return nil
}
