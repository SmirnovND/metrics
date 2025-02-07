package main

import (
	"fmt"
	config "github.com/SmirnovND/metrics/internal/pkg/config/agent"
	"github.com/SmirnovND/metrics/internal/usecase/agent"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
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

	cf := config.NewConfigCommand()
	go func() {
		agent.MetricsTracking(cf)
	}()
	http.ListenAndServe(addr, nil) // запускаем сервер
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
