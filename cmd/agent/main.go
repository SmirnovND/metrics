package main

import (
	"context"
	"fmt"
	config "github.com/SmirnovND/metrics/internal/pkg/config/agent"
	"github.com/SmirnovND/metrics/internal/usecase/agent"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

const (
	addr = ":8080" // адрес сервера
)

func main() {
	// Вывод информации о сборке
	printBuildInfo()

	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cf := config.NewConfigCommand()

	// Запускаем агент в фоне
	go agent.MetricsTracking(ctx, cf)

	// Настраиваем HTTP-сервер
	server := &http.Server{Addr: addr}

	// Канал для перехвата сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// Обрабатываем сигнал завершения
	go func() {
		<-sigChan
		fmt.Println("Получен сигнал завершения, агент корректно завершает работу...")

		// Останавливаем агент
		cancel()

		// Завершаем сервер с таймаутом
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("Ошибка при завершении сервера: %v\n", err)
		} else {
			fmt.Println("Сервер корректно завершен")
		}
	}()

	// Запускаем сервер (блокирующий вызов)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
	}

	fmt.Println("Приложение завершено")
}

// printBuildInfo выводит информацию о версии сборки
func printBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
