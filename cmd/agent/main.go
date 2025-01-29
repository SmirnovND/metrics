package main

import (
	config "github.com/SmirnovND/metrics/internal/pkg/config/agent"
	"github.com/SmirnovND/metrics/internal/usecase/agent"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
)

const (
	addr = ":8080" // адрес сервера
)

func main() {
	cf := config.NewConfigCommand()
	go func() {
		agent.MetricsTracking(cf)
	}()
	http.ListenAndServe(addr, nil) // запускаем сервер
}
