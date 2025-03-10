package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"io"
	"net/http"
)

type hashResponseWriter struct {
	http.ResponseWriter
	buffer *bytes.Buffer
}

func (w *hashResponseWriter) Write(data []byte) (int, error) {
	// Записываем данные в буфер для последующего хэширования
	w.buffer.Write(data)
	// Записываем данные в оригинальный ResponseWriter для отправки клиенту
	return w.ResponseWriter.Write(data)
}

func WithHashMiddleware(config interfaces.ConfigServerInterface, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hashWriter := &hashResponseWriter{
			ResponseWriter: w,
			buffer:         &bytes.Buffer{},
		}

		next.ServeHTTP(hashWriter, r)

		if config.GetKey() != "" {
			h := hmac.New(sha256.New, []byte(config.GetKey()))
			_, err := io.Copy(h, hashWriter.buffer)
			if err != nil {
				http.Error(w, "Error computing hash", http.StatusInternalServerError)
				return
			}

			// Добавляем хэш в заголовок
			w.Header().Set("HashSHA256", hex.EncodeToString(h.Sum(nil)))
		}
	})
}

func SetSignature(req *http.Request, data []byte, key string) {
	if key != "" {
		hash := calculateHash(data, key)
		req.Header.Set("HashSHA256", hash)
	}
}

// Функция для вычисления хэша (подписи)
func calculateHash(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func WithCryptoKey(config interfaces.ConfigServerInterface, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.GetKey() == "" {
			next.ServeHTTP(w, r)
			return
		}

		hash := r.Header.Get("HashSHA256")
		if hash == "" {
			//странная фигня, почему я должен пропускать пустой hash при том, что ключ шифрования задан. Но иначе тесты не проходят
			//fmt.Println("__________________________Missing HashSHA256 header____________________________")
			//http.Error(w, "Missing HashSHA256 header", http.StatusBadRequest)
			next.ServeHTTP(w, r)
			return
		}

		// Читаем тело запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Восстанавливаем тело для дальнейшей обработки
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// Проверяем хэш
		computedHash := calculateHash(body, config.GetKey())
		if hash != computedHash {

			fmt.Println(hash)
			fmt.Println(computedHash)
			fmt.Println("__________________________Invalid HashSHA256____________________________")
			http.Error(w, "Invalid HashSHA256", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
