package main

import (
	_ "github.com/SmirnovND/metrics/docs"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/middleware"
	"github.com/SmirnovND/metrics/internal/pkg/compressor"
	"github.com/SmirnovND/metrics/internal/pkg/container"
	"github.com/SmirnovND/metrics/internal/pkg/crypto"
	"github.com/SmirnovND/metrics/internal/pkg/loggeer"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/SmirnovND/metrics/internal/router"
	usecase "github.com/SmirnovND/metrics/internal/usecase/server"
	"github.com/jmoiron/sqlx"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	if err := Run(); err != nil {
		panic(err)
	}
}

func Run() error {

	diContainer := container.NewContainer(container.WithStartCollectionFunc(usecase.RestoreBackup))

	var cf interfaces.ConfigServer
	var storage *repo.MemStorage
	var db *sqlx.DB
	diContainer.Invoke(func(c interfaces.ConfigServer, s *repo.MemStorage, d *sqlx.DB) {
		cf = c
		storage = s
		db = d
	})

	defer usecase.Backup(cf, storage, db)
	stopCh := make(chan struct{})
	defer close(stopCh)

	usecase.TimedBackup(cf, storage, db, stopCh)
	return http.ListenAndServe(cf.GetFlagRunAddr(), middleware.ChainMiddleware(
		router.Handler(storage, db, cf.GetFlagRunAddr()),
		loggeer.WithLogging,
		compressor.WithDecompression,
		compressor.WithCompression,
		func(next http.Handler) http.Handler {
			return crypto.WithCryptoKey(cf, next)
		},
		func(next http.Handler) http.Handler {
			return crypto.WithHashMiddleware(cf, next)
		},
	))
}
