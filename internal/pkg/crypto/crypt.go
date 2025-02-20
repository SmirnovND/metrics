package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Загрузка приватного ключа с паролем
func loadPrivateKeyWithPassword(filename, password string) (*rsa.PrivateKey, error) {
	// Чтение закрытого ключа из файла
	keyFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open private key file: %v", err)
	}
	defer keyFile.Close()

	keyBytes, err := ioutil.ReadAll(keyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key file: %v", err)
	}

	// Декодирование PEM данных
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	// Расшифровка с паролем, если необходимо
	var privKey *rsa.PrivateKey
	if x509.IsEncryptedPEMBlock(block) {
		decryptedBytes, err := x509.DecryptPEMBlock(block, []byte(password))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt private key: %v", err)
		}
		block.Bytes = decryptedBytes
	}

	// Парсинг приватного ключа
	privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return privKey, nil
}

func decryptPrivateKey(encryptedData []byte, password []byte) (*rsa.PrivateKey, error) {
	// Расшифровываем с использованием пароля
	decryptedKey, err := x509.DecryptPEMBlock(&pem.Block{Bytes: encryptedData}, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt PEM block: %v", err)
	}

	// Парсим приватный ключ в формате PKCS#1
	privKey, err := x509.ParsePKCS1PrivateKey(decryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return privKey, nil
}

// WithDecryption - Middleware для расшифровки данных с логированием
func WithDecryption(key string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логируем начало обработки запроса
		log.Println("Starting decryption middleware")

		// Загружаем приватный ключ
		privKey, err := loadPrivateKeyWithPassword(".cert/encrypted_private_key.pem", "secretkey")
		if err != nil {
			log.Printf("Error loading private key: %v\n", err)
			http.Error(w, "Private key loading error", http.StatusInternalServerError)
			return
		}
		log.Println("Private key loaded successfully")

		// Чтение зашифрованных данных из тела запроса
		encryptedData, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v\n", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		log.Println("Encrypted data read successfully, length:", len(encryptedData))

		// Расшифровка данных
		decryptedData, err := DecryptWithRSA(privKey, encryptedData)
		if err != nil {
			log.Printf("Error during decryption: %v\n", err)
			http.Error(w, "Decryption error", http.StatusInternalServerError)
			return
		}
		log.Println("Data decrypted successfully, length:", len(decryptedData))

		// Возвращаем расшифрованные данные в запрос
		r.Body = io.NopCloser(bytes.NewReader(decryptedData))
		r.ContentLength = int64(len(decryptedData))

		// Логируем продолжение выполнения запроса
		log.Println("Passing control to the next handler")
		next.ServeHTTP(w, r)

		// Логируем успешную обработку запроса
		log.Println("Decryption middleware processed the request successfully")
	})
}

// Функция для расшифровки данных с использованием RSA приватного ключа
func DecryptWithRSA(privKey *rsa.PrivateKey, encryptedData []byte) ([]byte, error) {
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %v", err)
	}
	return decryptedData, nil
}

// Функция для шифрования данных с использованием публичного RSA ключа
func EncryptWithRSA(pubKey *rsa.PublicKey, data []byte) ([]byte, error) {
	// Шифруем данные с использованием публичного ключа
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %v", err)
	}
	return encryptedData, nil
}
