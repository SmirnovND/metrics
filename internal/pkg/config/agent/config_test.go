package agent

import (
	"flag"
	"os"
	"testing"
)

func init() {
	// Сброс флагов командной строки для тестов
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestNewConfigCommand(t *testing.T) {
	// Сброс флагов командной строки для тестов
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Устанавливаем переменные окружения для теста
	os.Setenv("REPORT_INTERVAL", "15")
	os.Setenv("POLL_INTERVAL", "5")
	os.Setenv("ADDRESS", "localhost:9090")
	os.Setenv("KEY", "test_key")
	os.Setenv("RATE_LIMIT", "3")

	// Вызовем NewConfigCommand, чтобы создать новый конфиг
	cf := NewConfigCommand()

	// Проверка значений флагов и переменных окружения
	if cf.GetReportInterval() != 15 {
		t.Errorf("Expected ReportInterval to be 15, but got %v", cf.GetReportInterval())
	}

	if cf.GetPollInterval() != 5 {
		t.Errorf("Expected PollInterval to be 5, but got %v", cf.GetPollInterval())
	}

	if cf.GetServerHost() != "http://localhost:9090" {
		t.Errorf("Expected ServerHost to be 'http://localhost:9090', but got %v", cf.GetServerHost())
	}

	if cf.GetKey() != "test_key" {
		t.Errorf("Expected Key to be 'test_key', but got %v", cf.GetKey())
	}

	if cf.GetRateLimit() != 3 {
		t.Errorf("Expected RateLimit to be 3, but got %v", cf.GetRateLimit())
	}
}
