package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"io"
	"net/http"
	"os"
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

// Генерация случайного симметричного ключа для AES
func generateAESKey() ([]byte, error) {
	key := make([]byte, 32) // 256 бит для AES
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %v", err)
	}
	return key, nil
}

// Шифрование данных с использованием AES
func encryptAES(plainText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Инициализация вектора (IV) для AES
	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	_, err = rand.Read(iv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	// Шифрование данных
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plainText)

	return ciphertext, nil
}

// Загрузка публичного ключа из файла в формате PKIX
func loadPublicKey() (*rsa.PublicKey, error) {
	// Чтение публичного ключа из файла
	pubKeyFile, err := os.Open(".cert/public_key.pem")
	if err != nil {
		return nil, fmt.Errorf("unable to open public key file: %v", err)
	}
	defer pubKeyFile.Close()

	pubKeyBytes, err := io.ReadAll(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key file: %v", err)
	}

	// Декодируем PEM данные
	block, _ := pem.Decode(pubKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	// Парсим публичный ключ с использованием PKCS1 формата
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	// Преобразуем в rsa.PublicKey, если это возможно
	return pubKey, nil
}

// Шифрование симметричного ключа с использованием RSA
func encryptRSA(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt with RSA: %v", err)
	}
	return encryptedData, nil
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

	//// Генерация симметричного ключа AES
	//aesKey, err := generateAESKey()
	//if err != nil {
	//	fmt.Println("Ошибка при генерации AES ключа:", err)
	//	return
	//}
	//
	//// Шифрование данных с использованием AES
	//encryptedData, err := encryptAES(jsonData, aesKey)
	//if err != nil {
	//	fmt.Println("Ошибка при шифровании данных:", err)
	//	return
	//}
	//
	//// Загрузка публичного ключа для шифрования AES ключа
	//pubKey, err := loadPublicKey()
	//if err != nil {
	//	fmt.Println("Ошибка при загрузке публичного ключа:", err)
	//	return
	//}

	//// Шифрование AES ключа с использованием RSA
	//encryptedAESKey, err := encryptRSA(pubKey, aesKey)
	//if err != nil {
	//	fmt.Println("Ошибка при шифровании AES ключа:", err)
	//	return
	//}
	//
	//// Кодирование зашифрованного AES ключа в Base64
	//encodedAESKey := base64.StdEncoding.EncodeToString(encryptedAESKey)

	// Создание HTTP-запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	//// Добавление зашифрованного и закодированного AES ключа в заголовок
	//req.Header.Set("X-Aes-Key", encodedAESKey)

	// Установка подписи для запроса
	//crypto.SetSignature(req, encryptedData, key)
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
