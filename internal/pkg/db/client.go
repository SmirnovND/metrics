package db

import (
	"fmt"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

const (
	maxIdleConnections = 10
	maxOpenConnections = 100
)

func ConfigureDB(db *sqlx.DB) {
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetMaxOpenConns(maxOpenConnections)
}

func NewDB(c interfaces.ConfigServerInterface) *sqlx.DB {
	dsn := c.GetDBDsn()
	if c.GetDBDsn() == "" {
		dsn = "invalid_dsn"
	}

	db, err := sqlx.Open(
		"postgres",
		dsn,
	)

	if err != nil {
		return db
	}

	ConfigureDB(db)
	fmt.Println("DB connection success!")

	err = db.Ping()
	if err == nil {
		// Получаем *sql.DB из *sqlx.DB
		dbBase := db.DB
		driver, err := postgres.WithInstance(dbBase, &postgres.Config{})
		if err != nil {
			log.Info().
				Err(err).
				Msg("Ошибка создания драйвера миграции ")
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://migrations", // Путь к миграциям
			"postgres",          // Имя базы данных
			driver,
		)
		if err != nil {
			log.Info().
				Err(err).
				Msg("Ошибка инициализации миграции ")
		}

		// Выполняем миграции
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Info().
				Err(err).
				Msg("Ошибка выполнения миграции ")
		}

		fmt.Println("Миграции применены успешно!")
	}

	return db
}
