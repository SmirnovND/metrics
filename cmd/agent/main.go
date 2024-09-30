package main

import (
	"fmt"
	"github.com/SmirnovND/metrics/internal/pkg/config"
	"os"
	"path/filepath"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	configPath := filepath.Join(cwd, "cmd", "agent", "config.yaml")
	//cf := config.NewConfig("./cmd/agent/config.yaml")
	cf := config.NewConfig(configPath)
	usecase.TrackingMetrics(cf)
	// Блокировка главной горутины
}
