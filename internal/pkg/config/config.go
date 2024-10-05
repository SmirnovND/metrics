package config

import (
	"flag"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"gopkg.in/yaml.v3"
	"log"
	"os"
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

func NewConfigYaml(patch string) (cf interfaces.ConfigAgent) {
	config := new(Config)
	config.LoadConfig(patch)
	return config
}

func (c *Config) LoadConfig(patch string) {
	file, err := os.Open(patch)
	if err != nil {
		log.Fatal("ReadConfigFile: ", err)
	}

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatal("DecodeConfigFile: ", err)
	}
}

func NewConfigCommand() (cf interfaces.ConfigAgent) {
	config := new(Config)
	flag.StringVar(&config.ServerHost, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&config.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&config.PollInterval, "p", 2, "poll interval")
	return config
}
