package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Расшифровка AES-ключа с помощью RSA
func decryptAESKey(encryptedKey string, privateKey *rsa.PrivateKey) ([]byte, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedKey)
}

// Расшифровка JSON с помощью AES
func decryptAES(encryptedData string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext, nil
}

// WithDecryption — middleware для расшифровки данных запроса
func WithDecryption(privateKey *rsa.PrivateKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Читаем зашифрованные данные из тела запроса
			encryptedBody, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Println("Ошибка чтения тела запроса")
				http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
				return
			}
			r.Body.Close() // Закрываем исходное тело

			// Разбираем JSON: {"aes_key": "...", "data": "..."}
			payload, err := parseEncryptedPayload(encryptedBody)
			if err != nil {
				fmt.Println("Ошибка парсинга зашифрованного JSON")
				http.Error(w, "Ошибка парсинга зашифрованного JSON", http.StatusBadRequest)
				return
			}

			// Расшифровываем AES-ключ с помощью приватного RSA-ключа
			aesKey, err := decryptAESKey(payload.AESKey, privateKey)
			if err != nil {
				fmt.Println("Ошибка расшифровки AES-ключа", err)
				http.Error(w, "Ошибка расшифровки AES-ключа", http.StatusInternalServerError)
				return
			}

			// Расшифровываем данные с помощью AES
			decryptedData, err := decryptAES(payload.Data, aesKey)
			if err != nil {
				fmt.Println("Ошибка расшифровки данных")
				http.Error(w, "Ошибка расшифровки данных", http.StatusInternalServerError)
				return
			}

			// Подменяем тело запроса расшифрованными данными
			r.Body = io.NopCloser(bytes.NewReader(decryptedData))
			r.ContentLength = int64(len(decryptedData))
			fmt.Println("----------------------4444")
			// Передаем управление следующему middleware
			next.ServeHTTP(w, r)
		})
	}
}

// parseEncryptedPayload парсит входные зашифрованные данные
func parseEncryptedPayload(data []byte) (struct {
	AESKey string `json:"aes_key"`
	Data   string `json:"data"`
}, error) {
	var payload struct {
		AESKey string `json:"aes_key"`
		Data   string `json:"data"`
	}
	err := json.Unmarshal(data, &payload)
	return payload, err
}

// LoadPrivateKey загружает приватный RSA-ключ из файла
func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	// Читаем содержимое файла
	keyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Разбираем PEM-блок
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("не удалось декодировать PEM-блок")
	}

	// Парсим приватный ключ в формате PKCS#1 или PKCS#8
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		privateKeyPKCS8, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.New("не удалось разобрать приватный ключ: " + err.Error())
		}
		privateKey, ok := privateKeyPKCS8.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("ключ не является RSA-ключом")
		}
		return privateKey, nil
	}

	return privateKey, nil
}

// Генерация случайного AES-ключа (32 байта для AES-256)
func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	return key, err
}

// AES-шифрование данных
func EncryptAES(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Генерируем IV (инициализационный вектор)
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}

	// Шифрование в режиме CBC
	mode := cipher.NewCBCEncrypter(block, iv)
	paddedData := pkcs7Padding(data, aes.BlockSize)
	ciphertext := make([]byte, len(paddedData))
	mode.CryptBlocks(ciphertext, paddedData)

	// Объединяем IV + зашифрованный текст
	encryptedData := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// PKCS7 padding (дополнение до размера блока)
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// Функция загрузки публичного RSA-ключа
func LoadPublicKey(filename string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("не удалось декодировать PEM")
	}

	// Используем ParsePKIXPublicKey вместо ParsePKCS1PublicKey
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Приводим к типу *rsa.PublicKey
	pubKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("неверный тип ключа, ожидался *rsa.PublicKey")
	}

	return pubKey, nil
}

// RSA-шифрование AES-ключа
func EncryptAESKey(aesKey []byte, publicKey *rsa.PublicKey) (string, error) {
	encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, aesKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedKey), nil
}
