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

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		config.ServerHost = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		reportInterval, err := strconv.Atoi(envReportInterval)
		if err == nil {
			config.ReportInterval = reportInterval
		}
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		envPollInterval, err := strconv.Atoi(envPollInterval)
		if err == nil {
			config.PollInterval = envPollInterval
		}
	}

	config.ServerHost = fmt.Sprintf("http://%s", config.ServerHost)
	return config
}
