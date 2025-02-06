package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"os"
)

type ServiceBackup struct {
	storage interfaces.MemStorageInterface
	cf      interfaces.ConfigServer
	db      *sqlx.DB
}

func NewServiceBackup(storage interfaces.MemStorageInterface, cf interfaces.ConfigServer, db *sqlx.DB) *ServiceBackup {
	return &ServiceBackup{
		storage: storage,
		cf:      cf,
		db:      db,
	}
}

func (s *ServiceBackup) Backup() {
	err := s.db.Ping()
	if err != nil {
		log.Info().
			Err(err).
			Msg("Ошибка соединения с базой ")
		s.backupToFile()
		return
	} else {
		s.backupToDB()
	}
}

func (s *ServiceBackup) backupToFile() {
	file, err := os.Create(s.cf.GetFileStoragePath())
	if err != nil {
		log.Info().
			Err(err).
			Msg("Ошибка backupToFile ")
		return
	}
	defer file.Close()
	s.storage.ExecuteWithLock(func(collection map[string]domain.MetricInterface) {
		encoder := json.NewEncoder(file)
		err = encoder.Encode(collection)
		if err != nil {
			log.Info().
				Err(err).
				Msg("Ошибка backupToFile ")
			return
		}
	})
}

func (s *ServiceBackup) backupToDB() {
	tx, err := s.db.Begin()
	if err != nil {
		log.Info().
			Err(err).
			Msg("Ошибка старта транзакции ")
		return
	}

	_, err = tx.Exec("DELETE FROM metric")
	if err != nil {
		tx.Rollback()
		log.Info().
			Err(err).
			Msg("Ошибка очистки таблицы ")
		return
	}

	s.storage.ExecuteWithLock(func(collection map[string]domain.MetricInterface) {
		stmt, err := tx.Prepare("INSERT INTO metric (id, type, value, delta) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET type = EXCLUDED.type, value = EXCLUDED.value, delta = EXCLUDED.delta")
		if err != nil {
			log.Info().
				Err(err).
				Msg("Ошибка подготовки запроса ")
			tx.Rollback()
			return
		}
		defer stmt.Close()

		for _, metricInterface := range collection {
			metric := metricInterface.(*domain.Metric)
			var delta sql.NullInt64
			var value sql.NullFloat64

			if metric.Delta != 0 {
				delta = sql.NullInt64{Int64: metric.Delta, Valid: true}
			}
			if metric.Value != 0.0 {
				value = sql.NullFloat64{Float64: metric.Value, Valid: true}
			}

			_, errExec := stmt.Exec(metric.ID, metric.MType, value, delta)
			if errExec != nil {
				log.Info().
					Err(errExec).
					Msg("Ошибка вставки ")
				tx.Rollback()
				return
			}
		}

		errC := tx.Commit()
		if errC != nil {
			log.Info().
				Err(errC).
				Msg("Ошибка комита транщакции ")
		} else {
			log.Info().
				Err(errC).
				Msg("Метрики сохранены ")
		}
	})
}

func RestoreMetric(cf interfaces.ConfigServer, db *sqlx.DB) map[string]domain.MetricInterface {
	err := db.Ping()
	if err != nil {
		return RestoreFromFile(cf.GetFileStoragePath(), cf.IsRestore())
	} else {
		return RestoreFromDB(db)
	}
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

func RestoreFromDB(db *sqlx.DB) map[string]domain.MetricInterface {
	rows, err := db.Query("SELECT id, type, value, delta FROM metric")
	if err != nil {
		fmt.Println("Error querying database:", err)
		return make(map[string]domain.MetricInterface)
	}
	defer rows.Close()

	collection := make(map[string]domain.MetricInterface)

	for rows.Next() {
		var (
			id    string
			mtype string
			delta sql.NullInt64
			value sql.NullFloat64
		)

		if err := rows.Scan(&id, &mtype, &value, &delta); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}

		metric := &domain.Metric{
			ID:    id,
			MType: mtype,
		}

		if delta.Valid {
			metric.Delta = delta.Int64
		}
		if value.Valid {
			metric.Value = value.Float64
		}

		collection[id] = metric
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error with rows:", err)
	}

	return collection
}
