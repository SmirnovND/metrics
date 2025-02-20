package server_test

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/mocks"
	"github.com/SmirnovND/metrics/internal/repo"
	"github.com/SmirnovND/metrics/internal/services/server"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestBackupToFile(t *testing.T) {
	// Создаем моки
	mockStorage := new(mocks.MemStorageInterface)
	mockConfig := new(mocks.ConfigServer)

	// Создаем временный файл для проверки записи данных
	tmpFile, err := os.CreateTemp("", "backup_test.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Ожидание вызова GetFileStoragePath() и возврат пути файла
	mockConfig.On("GetFileStoragePath").Return(tmpFile.Name())

	// Эмулируем ExecuteWithLock, передавая внутрь функцию
	mockStorage.On("ExecuteWithLock", mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(0).(func(map[string]domain.MetricInterface))
		testData := map[string]domain.MetricInterface{
			"test_metric": nil, // Добавляем тестовые данные
		}
		fn(testData) // Вызываем переданную функцию с тестовыми данными
	}).Once()

	// Создаем сервис
	service := server.NewServiceBackup(mockStorage, mockConfig, nil)

	// Вызываем BackupToFile()
	service.BackupToFile()

	// Читаем записанный файл
	content, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)
	require.Contains(t, string(content), "test_metric") // Проверяем, что данные записались
}

func TestBackupToFile_FileCreateError(t *testing.T) {
	mockStorage := new(MockMemStorage)
	mockConfig := new(MockConfigServer)

	// Неверный путь, эмулируем ошибку создания файла
	mockConfig.On("GetFileStoragePath").Return("/invalid_path/backup.json")

	service := server.NewServiceBackup(mockStorage, mockConfig, nil)
	service.BackupToFile()

	// Проверяем, что `ExecuteWithLock` не вызывался, так как файл не создался
	mockStorage.AssertNotCalled(t, "ExecuteWithLock", mock.Anything)
	mockConfig.AssertExpectations(t)
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

// MockConfigServer - мок для ConfigServer
type MockConfigServer struct {
	mock.Mock
}

func (m *MockConfigServer) GetFileStoragePath() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigServer) GetStoreInterval() time.Duration {
	return 0 // Если в тесте не используется, просто возвращаем 0
}

func (m *MockConfigServer) IsRestore() bool {
	return false // Аналогично
}

func (m *MockConfigServer) GetFlagRunAddr() string {
	return "" // Аналогично
}

func (m *MockConfigServer) GetDBDsn() string {
	return "" // Заглушка
}

func (m *MockConfigServer) GetKey() string {
	return "" // Заглушка
}
