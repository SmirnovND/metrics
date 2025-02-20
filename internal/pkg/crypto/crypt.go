package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// LoadPublicKey загружает публичный ключ из файла
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("неверный формат публичного ключа")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("неверный тип публичного ключа")
	}

	return rsaPub, nil
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("неверный формат приватного ключа")
	}

	// Попробуем сначала как PKCS#1
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return priv, nil
	}

	// Если ошибка, попробуем как PKCS#8
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("не удалось распарсить приватный ключ в формате PKCS#1 или PKCS#8")
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("загруженный ключ не является RSA-ключом")
	}

	return rsaKey, nil
}

// WithDecryption Middleware для расшифровки входящих данных
func WithDecryption(path string, next http.Handler) http.Handler {
	// Загружаем приватный ключ
	privateKey, err := LoadPrivateKey(path)
	if err != nil {
		// Логируем ошибку загрузки ключа и возвращаем следующий обработчик
		fmt.Printf("Ошибка загрузки приватного ключа: %v\n", err)
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Чтение зашифрованного тела запроса
		encryptedData, err := io.ReadAll(r.Body)
		if err != nil {
			// Логируем ошибку чтения тела запроса
			http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Расшифровка данных
		decryptedData, err := Decrypt(encryptedData, privateKey)
		if err != nil {
			// Логируем ошибку расшифровки данных
			http.Error(w, "Ошибка расшифровки данных", http.StatusUnauthorized)
			fmt.Println(err)
			return
		}

		// Создаем новый запрос с расшифрованным телом
		r.Body = io.NopCloser(bytes.NewReader(decryptedData))
		r.ContentLength = int64(len(decryptedData))

		// Продолжаем выполнение следующего обработчика
		next.ServeHTTP(w, r)
	})
}

// EncryptRSA Функция для шифрования данных с использованием RSA
func EncryptRSA(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	// Используем PKCS1v15 для шифрования данных
	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при шифровании RSA: %v", err)
	}
	return encryptedData, nil
}

// EncryptAES шифрует данные с использованием AES и возвращает только зашифрованные данные и ошибку
func EncryptAES(data []byte, key []byte) ([]byte, error) {
	// Генерация случайного вектора инициализации (IV) для шифрования
	iv := make([]byte, aes.BlockSize)
	_, err := rand.Read(iv)
	if err != nil {
		return nil, fmt.Errorf("не удалось сгенерировать IV: %v", err)
	}

	// Создание AES-блока шифрования
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания шифратора AES: %v", err)
	}

	// Подготовка шифрования в режиме CBC
	stream := cipher.NewCBCEncrypter(block, iv)

	// Дополним данные до размера, кратного размеру блока
	padding := aes.BlockSize - len(data)%aes.BlockSize
	paddedData := append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)

	// Шифруем данные
	ciphertext := make([]byte, len(paddedData))
	stream.CryptBlocks(ciphertext, paddedData)

	// Возвращаем зашифрованные данные, включая IV
	// Мы возвращаем только зашифрованные данные и IV как одно целое
	return append(iv, ciphertext...), nil
}

// GenerateSignature генерирует HMAC-подпись на основе данных и секретного ключа
func GenerateSignature(data []byte, secretKey string) string {
	// Создаем HMAC с использованием SHA-256
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Decrypt расшифровывает данные с помощью приватного ключа
func Decrypt(data []byte, priv *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, data, nil)
}
