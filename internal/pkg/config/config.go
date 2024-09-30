package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	ReportInterval int    `yaml:"reportInterval"`
	PollInterval   int    `yaml:"pollInterval"`
	ServerHost     string `yaml:"serverHost"`
}

func NewConfig(patch string) (cf Config) {
	cf.LoadConfig(patch)
	return cf
}

func (cg *Config) LoadConfig(patch string) {
	file, err := os.Open(patch)
	if err != nil {
		log.Fatal("ReadConfigFile: ", err)
	}

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cg)
	if err != nil {
		log.Fatal("DecodeConfigFile: ", err)
	}
}
