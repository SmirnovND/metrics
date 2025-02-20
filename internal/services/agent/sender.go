package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	cryptoPkg "github.com/SmirnovND/metrics/internal/pkg/crypto"
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

func SendJSON(m *domain.Metrics, serverHost string, key string, cryptoKeyPath string) {
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

	// Загрузка публичного ключа
	publicKey, err := cryptoPkg.LoadPublicKey(cryptoKeyPath)
	if err != nil {
		fmt.Println("Ошибка загрузки публичного ключа:", err)
		return
	}

	// Генерация симметричного ключа AES для шифрования данных
	aesKey := make([]byte, 32)
	_, err = rand.Read(aesKey)
	if err != nil {
		fmt.Println("Ошибка при генерации AES-ключа:", err)
		return
	}

	// Шифрование данных с использованием AES
	encryptedData, err := cryptoPkg.EncryptAES(jsonData, aesKey)
	if err != nil {
		fmt.Println("Ошибка при шифровании данных AES:", err)
		return
	}

	// Шифрование AES-ключа с использованием RSA
	encryptedKey, err := cryptoPkg.EncryptRSA(publicKey, aesKey)
	if err != nil {
		fmt.Println("Ошибка при шифровании ключа RSA:", err)
		return
	}

	// Создание структуры для отправки
	payload := struct {
		Key       string `json:"key"`
		Nonce     string `json:"nonce"`
		Data      string `json:"data"`
		Signature string `json:"signature"`
	}{
		Key:       base64.StdEncoding.EncodeToString(encryptedKey),
		Nonce:     "", // Удалил nonce, так как его нет
		Data:      base64.StdEncoding.EncodeToString(encryptedData),
		Signature: cryptoPkg.GenerateSignature(jsonData, key),
	}

	// Сериализация в JSON
	payloadData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Ошибка при сериализации payload:", err)
		return
	}

	// Создание HTTP-запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadData))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
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
