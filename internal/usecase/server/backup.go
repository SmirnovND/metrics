package server

import (
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/SmirnovND/metrics/internal/services/server"
	"github.com/jmoiron/sqlx"
	"time"
)

func TimedBackup(cf interfaces.ConfigServer, storage *repo.MemStorage, db *sqlx.DB, stopCh <-chan struct{}) {
	backupTicker := time.NewTicker(cf.GetStoreInterval())
	service := server.NewServiceBackup(storage, cf, db)

	go func() {
		defer backupTicker.Stop()
		for {
			select {
			case <-backupTicker.C:
				service.Backup()
				fmt.Println("Выполняется резервное копирование...")
			case <-stopCh:
				return
			}
		}
	}()
}

func Backup(cf interfaces.ConfigServer, storage *repo.MemStorage, db *sqlx.DB) {
	service := server.NewServiceBackup(storage, cf, db)
	service.Backup()
}

func RestoreBackup(cf interfaces.ConfigServer, db *sqlx.DB) map[string]domain.MetricInterface {
	return server.RestoreMetric(cf, db)
}
