package server_test

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/SmirnovND/metrics/internal/services/server"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestServiceBackup_Backup(t *testing.T) {
	// Создание мок-соединения с базой данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Преобразуем стандартное соединение *sql.DB в *sqlx.DB
	sqlxDB := sqlx.NewDb(db, "postgres")

	// Создаем мок-конфиг
	cf := &mockConfigServer{}

	// Создаем экземпляр ServiceBackup
	storage := &repo.MemStorage{}
	service := server.NewServiceBackup(storage, cf, sqlxDB)

	// Вызов метода Backup
	service.Backup()

	// Проверяем, что ожидания для mock были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания были выполнены: %v", err)
	}
}

func TestRestoreMetric(t *testing.T) {
	t.Run("Restore from file", func(t *testing.T) {
		// Мокируем конфигурацию
		cf := &mockConfigServer{}

		// Подготавливаем файл для чтения
		metrics := map[string]*domain.Metric{
			"metric1": {ID: "metric1", MType: "counter", Value: 10},
		}
		file, err := os.Create(cf.GetFileStoragePath())
		if err != nil {
			t.Fatal("Не удалось создать файл для теста:", err)
		}
		defer os.Remove(file.Name()) // Удаляем файл после теста

		// Записываем в файл
		encoder := json.NewEncoder(file)
		if err := encoder.Encode(metrics); err != nil {
			t.Fatal("Ошибка записи в файл:", err)
		}

		// Создание мок-соединения с базой данных
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		db.Close()

		// Преобразуем стандартное соединение *sql.DB в *sqlx.DB
		sqlxDB := sqlx.NewDb(db, "postgres")

		// Вызываем функцию RestoreMetric
		result := server.RestoreMetric(cf, sqlxDB)

		// Проверяем, что результат верен
		assert.Len(t, result, 1)
		assert.Equal(t, "metric1", result["metric1"].(*domain.Metric).ID)
		assert.Equal(t, 10.0, result["metric1"].(*domain.Metric).Value)
	})

	t.Run("Restore from DB", func(t *testing.T) {
		// Мокируем соединение с базой данных
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal("Не удалось создать мок-соединение с базой данных:", err)
		}
		defer db.Close()

		// Мокаем конфигурацию
		cf := &mockConfigServer{}

		// Мокаем запрос к базе данных
		rows := mock.NewRows([]string{"id", "type", "value", "delta"}).
			AddRow("metric1", "counter", 10.0, nil)
		mock.ExpectQuery("SELECT id, type, value, delta FROM metric").
			WillReturnRows(rows)

		// Вызываем функцию RestoreMetric
		result := server.RestoreMetric(cf, sqlx.NewDb(db, "postgres"))

		// Проверяем, что результат верен
		assert.Len(t, result, 1)
		assert.Equal(t, "metric1", result["metric1"].(*domain.Metric).ID)
		assert.Equal(t, 10.0, result["metric1"].(*domain.Metric).Value)

		// Проверяем, что ожидания для mock выполнены
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("Не все ожидания были выполнены: %v", err)
		}
	})
}

type mockConfigServer struct{}

func (m *mockConfigServer) GetFileStoragePath() string {
	return "/tmp/backup.json"
}

func (m *mockConfigServer) IsRestore() bool {
	return true
}

func (m *mockConfigServer) GetStoreInterval() time.Duration {
	return time.Second * 10
}

func (m *mockConfigServer) GetFlagRunAddr() string {
	return "localhost:8080"
}

func (m *mockConfigServer) GetDBDsn() string {
	return "user:password@tcp(localhost:3306)/dbname"
}

func (m *mockConfigServer) GetKey() string {
	return "some-key"
}
