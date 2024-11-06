package server

import (
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/repo"
	"os"
)

type ServiceBackup struct {
	storage *repo.MemStorage
	cf      interfaces.ConfigServer
}

func NewServiceBackup(storage *repo.MemStorage, cf interfaces.ConfigServer) *ServiceBackup {
	return &ServiceBackup{
		storage: storage,
		cf:      cf,
	}
}

func (s *ServiceBackup) Backup() {
	file, err := os.Create(s.cf.GetFileStoragePath())
	if err != nil {
		fmt.Println("Error Backup:", err)
		return
	}
	defer file.Close()
	s.storage.ExecuteWithLock(func(collection map[string]domain.MetricInterface) {
		encoder := json.NewEncoder(file)
		err = encoder.Encode(collection)
		if err != nil {
			fmt.Println("Error Backup:", err)
			return
		}
	})
}

func RestoreFromFile(filename string, restore bool) map[string]domain.MetricInterface {
	newCollection := make(map[string]*domain.Metric)

	if !restore {
		return make(map[string]domain.MetricInterface)
	}

	file, err := os.Open(filename)
	if err != nil {
		return make(map[string]domain.MetricInterface)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&newCollection); err != nil {
		fmt.Printf("Ошибка декодирования JSON: %v\n", err)
		return make(map[string]domain.MetricInterface)
	}

	result := make(map[string]domain.MetricInterface, len(newCollection))
	for key, metric := range newCollection {
		result[key] = metric
	}

	return result
}
