package server

import (
	"flag"
	"os"
	"testing"
	"time"
)

func init() {
	// Сброс флагов командной строки для тестов
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestNewConfigCommand(t *testing.T) {
	// Сброс флагов командной строки для тестов
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Устанавливаем переменные окружения для теста
	os.Setenv("STORE_INTERVAL", "600")
	os.Setenv("FILE_STORAGE_PATH", "./test_storage.json")
	os.Setenv("DATABASE_DSN", "test_db_dsn")
	os.Setenv("RESTORE", "false")
	os.Setenv("ADDRESS", "localhost:9090")
	os.Setenv("KEY", "test_key")

	// Вызовем NewConfigCommand, чтобы создать новый конфиг
	cf := NewConfigCommand()

	// Проверка значений флагов и переменных окружения
	if cf.GetStoreInterval() != 600*time.Second {
		t.Errorf("Expected StoreInterval to be %v, but got %v", 600*time.Second, cf.GetStoreInterval())
	}

	if cf.GetFileStoragePath() != "./test_storage.json" {
		t.Errorf("Expected FileStoragePath to be './test_storage.json', but got %v", cf.GetFileStoragePath())
	}

	if cf.GetDBDsn() != "test_db_dsn" {
		t.Errorf("Expected DBDsn to be 'test_db_dsn', but got %v", cf.GetDBDsn())
	}

	if cf.IsRestore() != false {
		t.Errorf("Expected Restore to be false, but got %v", cf.IsRestore())
	}

	if cf.GetFlagRunAddr() != "localhost:9090" {
		t.Errorf("Expected FlagRunAddr to be 'localhost:9090', but got %v", cf.GetFlagRunAddr())
	}

	if cf.GetKey() != "test_key" {
		t.Errorf("Expected Key to be 'test_key', but got %v", cf.GetKey())
	}
}
