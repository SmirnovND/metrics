package server

import (
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/services/server"
	"github.com/jmoiron/sqlx"
	"time"
)

// TimedBackup выполняет резервное копирование метрик через заданные интервалы времени.
// Процесс запускается в отдельной горутине и продолжает работу до получения сигнала остановки.
func TimedBackup(cf interfaces.ConfigServerInterface, storage interfaces.MemStorageInterface, db *sqlx.DB, stopCh <-chan struct{}) {
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

// Backup выполняет разовое резервное копирование метрик.
func Backup(cf interfaces.ConfigServerInterface, storage interfaces.MemStorageInterface, db *sqlx.DB) {
	service := server.NewServiceBackup(storage, cf, db)
	service.Backup()
}

// RestoreBackup восстанавливает метрики из резервной копии.
func RestoreBackup(cf interfaces.ConfigServerInterface, db *sqlx.DB) map[string]domain.MetricInterface {
	return server.RestoreMetric(cf, db)
}
