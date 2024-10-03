package main

import (
	"fmt"
	"github.com/SmirnovND/metrics/internal/pkg/config"
	"github.com/SmirnovND/metrics/internal/services/tracking"
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
	cf := config.NewConfig(configPath)
	tracking.TrackingMetrics(cf)
}
