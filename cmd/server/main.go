package main

import (
	"flag"
	"github.com/SmirnovND/metrics/internal/middleware"
	"github.com/SmirnovND/metrics/internal/pkg/compressor"
	"github.com/SmirnovND/metrics/internal/pkg/loggeer"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/SmirnovND/metrics/internal/router"
	"net/http"
	"os"
)

func main() {
	if err := Run(); err != nil {
		panic(err)
	}
}

func Run() error {
	storage := repo.NewMetricRepo()

	var flagRunAddr string
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	return http.ListenAndServe(flagRunAddr, middleware.ChainMiddleware(
		router.Handler(storage),
		loggeer.WithLogging,
		compressor.WithDecompression,
		compressor.WithCompression,
	))
}
