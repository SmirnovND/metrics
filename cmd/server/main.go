package main

import (
	"github.com/SmirnovND/metrics/internal/middleware"
	"github.com/SmirnovND/metrics/internal/pkg/compressor"
	config "github.com/SmirnovND/metrics/internal/pkg/config/server"
	"github.com/SmirnovND/metrics/internal/pkg/loggeer"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/SmirnovND/metrics/internal/router"
	usecase "github.com/SmirnovND/metrics/internal/usecase/server"
	"net/http"
)

func main() {
	if err := Run(); err != nil {
		panic(err)
	}
}

func Run() error {
	cf := config.NewConfigCommand()
	storage := repo.NewMetricRepo(usecase.RestoreBackup(cf))
	defer usecase.Backup(cf, storage)
	stopCh := make(chan struct{})
	defer close(stopCh)

	usecase.TimedBackup(cf, storage, stopCh)
	return http.ListenAndServe(cf.GetFlagRunAddr(), middleware.ChainMiddleware(
		router.Handler(storage),
		loggeer.WithLogging,
		compressor.WithDecompression,
		compressor.WithCompression,
	))
}
