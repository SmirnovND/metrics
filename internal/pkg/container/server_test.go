package container

import (
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/jmoiron/sqlx"
	"testing"
)

// Моковая функция для инициализации коллекции метрик
func mockStartCollectionFunc(cf interfaces.ConfigServer, db *sqlx.DB) map[string]domain.MetricInterface {
	metric := &domain.Gauge{}
	value := 42.0
	metric.SetName("test_metric").SetValue(&value)
	return map[string]domain.MetricInterface{
		"test_metric": metric,
	}
}

func TestNewContainer(t *testing.T) {
	// Создаем контейнер с моковой функцией
	c := NewContainer(WithStartCollectionFunc(mockStartCollectionFunc))

	if c == nil {
		t.Fatal("Container не должен быть nil")
	}

	if c.container == nil {
		t.Fatal("Container должен содержать dig.Container")
	}

	if c.startCollectionFunc == nil {
		t.Fatal("startCollectionFunc не должна быть nil")
	}
}

func TestProvideDependencies(t *testing.T) {
	// Создаем контейнер с моковой функцией
	c := NewContainer(WithStartCollectionFunc(mockStartCollectionFunc))

	err := c.Invoke(func(cfg interfaces.ConfigServer, database *sqlx.DB, storage *repo.MemStorage) {
		if cfg == nil {
			t.Error("ConfigServer не должен быть nil")
		}
		if database == nil {
			t.Error("DB не должна быть nil")
		}
		if storage == nil {
			t.Error("MemStorage не должен быть nil")
		}
	})

	if err != nil {
		t.Fatalf("Ошибка при инжекте зависимостей: %v", err)
	}
}

func TestInvoke_Error(t *testing.T) {
	c := NewContainer(WithStartCollectionFunc(mockStartCollectionFunc))

	// Попытка инжектировать несуществующую зависимость
	err := c.Invoke(func(nonexistentDep string) {})

	if err == nil {
		t.Fatal("Ожидалась ошибка из-за отсутствующей зависимости")
	}
}
