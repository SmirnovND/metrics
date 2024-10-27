package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// gzipResponseWriter оборачивает стандартный ResponseWriter для добавления gzip-сжатия
type gzipResponseWriter struct {
	http.ResponseWriter
	*gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Header возвращает заголовки ответа
func (w gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// WithCompression добавляет сжатие ответов
func WithCompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Проверяем заголовок Content-Type и устанавливаем его, если он не был задан
		contentType := w.Header().Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html") {
			w.Header().Set("Content-Type", contentType)

			gz := gzip.NewWriter(w)
			defer gz.Close()

			gzResponseWriter := gzipResponseWriter{Writer: gz, ResponseWriter: w}
			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gzResponseWriter, r)
		} else {
			next.ServeHTTP(w, r)
			return
		}
	})
}

func WithDecompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Ошибка при распаковке gzip", http.StatusBadRequest)
				return
			}
			defer gz.Close()

			r.Body = io.NopCloser(gz)
		}

		next.ServeHTTP(w, r)
	})
}
