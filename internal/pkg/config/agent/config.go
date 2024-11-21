package agent

import (
	"flag"
	"fmt"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"os"
	"strconv"
)

const (
	reportInterval = 10
	pollInterval   = 2
)

type Config struct {
	ReportInterval int    `yaml:"reportInterval"`
	PollInterval   int    `yaml:"pollInterval"`
	ServerHost     string `yaml:"serverHost"`
	Key            string `yaml:"key"`
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
func (c *Config) GetKey() string {
	return c.Key
}

func NewConfigCommand() (cf interfaces.ConfigAgent) {
	config := new(Config)

	flag.StringVar(&config.ServerHost, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&config.ReportInterval, "r", reportInterval, "report interval")
	flag.IntVar(&config.PollInterval, "p", pollInterval, "poll interval")
	flag.StringVar(&config.Key, "k", "", "key")

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

	if envKey := os.Getenv("KEY"); envKey != "" {
		config.Key = envKey
	}

	config.ServerHost = fmt.Sprintf("http://%s", config.ServerHost)
	return config
}
