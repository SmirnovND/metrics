package container

import (
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
	config "github.com/SmirnovND/metrics/internal/pkg/config/server"
	"github.com/SmirnovND/metrics/internal/pkg/db"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/dig"
)

// Container - структура контейнера, обертывающая dig-контейнер
type Container struct {
	container           *dig.Container
	startCollectionFunc func(cf interfaces.ConfigServer, db *sqlx.DB) map[string]domain.MetricInterface
}

// NewContainer - создаёт новый DI-контейнер и регистрирует зависимости
func NewContainer(startCollectionFunc func(cf interfaces.ConfigServer, db *sqlx.DB) map[string]domain.MetricInterface) *Container {
	c := &Container{
		container:           dig.New(),
		startCollectionFunc: startCollectionFunc,
	}

	c.provideDependencies()
	return c
}

// provideDependencies - функция, регистрирующая зависимости
func (c *Container) provideDependencies() {
	// Регистрируем конфигурацию
	c.container.Provide(config.NewConfigCommand)

	// Регистрируем db
	c.container.Provide(db.NewDB)

	// Регистрируем репозиторий, передав конфигурацию
	c.container.Provide(func(cf interfaces.ConfigServer, db *sqlx.DB) *repo.MemStorage {
		return repo.NewMetricRepo(c.startCollectionFunc(cf, db))
	})
}

// Invoke - функция для вызова и инжекта зависимостей
func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}
