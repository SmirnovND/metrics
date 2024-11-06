package server

import (
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/SmirnovND/metrics/internal/services/server"
	"time"
)

func TimedBackup(cf interfaces.ConfigServer, storage *repo.MemStorage, stopCh <-chan struct{}) {
	backupTicker := time.NewTicker(cf.GetStoreInterval())
	service := server.NewServiceBackup(storage, cf)

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

func Backup(cf interfaces.ConfigServer, storage *repo.MemStorage) {
	service := server.NewServiceBackup(storage, cf)
	service.Backup()
}

func RestoreBackup(cf interfaces.ConfigServer) map[string]domain.MetricInterface {
	return server.RestoreFromFile(cf.GetFileStoragePath(), cf.IsRestore())
}
