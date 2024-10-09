package agent

import (
	"flag"
	"fmt"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"os"
	"strconv"
)

type Config struct {
	ReportInterval int    `yaml:"reportInterval"`
	PollInterval   int    `yaml:"pollInterval"`
	ServerHost     string `yaml:"serverHost"`
}

func (c *Config) GetReportInterval() int {
	return c.ReportInterval
}

func (c *Config) GetPollInterval() int {
	return c.PollInterval
}

func (c *Config) GetServerHost() string {
	return c.ServerHost
}

func NewConfigCommand() (cf interfaces.ConfigAgent) {
	config := new(Config)

	flag.StringVar(&config.ServerHost, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&config.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&config.PollInterval, "p", 2, "poll interval")

	flag.Parse()

	// Использование
	config.ServerHost = getEnvOrDefault("ADDRESS", "localhost:8080")

	envReportInterval := getEnvOrDefault("REPORT_INTERVAL", "10")
	reportInterval, err := strconv.Atoi(envReportInterval)
	if err == nil {
		config.ReportInterval = reportInterval
	}

	envPollInterval := getEnvOrDefault("POLL_INTERVAL", "2")
	pollInterval, err := strconv.Atoi(envPollInterval)
	if err == nil {
		config.PollInterval = pollInterval
	}

	config.ServerHost = fmt.Sprintf("http://%s", config.ServerHost)
	return config
}

func getEnvOrDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
