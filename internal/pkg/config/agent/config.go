package agent

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"io"
	"os"
	"strconv"
)

const (
	reportInterval = 10
	pollInterval   = 2
	rateLimit      = 1
)

type Config struct {
	ReportInterval int    `yaml:"reportInterval" json:"reportInterval"`
	PollInterval   int    `yaml:"pollInterval" json:"pollInterval"`
	ServerHost     string `yaml:"serverHost" json:"serverHost"`
	GRPCServerHost string `yaml:"GRPCServerHost" json:"GRPCServerHost"`
	Key            string `yaml:"key" json:"key"`
	RateLimit      int    `yaml:"rateLimit" json:"rateLimit"`
	CryptoKey      string `yaml:"cryptoKey" json:"cryptoKey"`
	UseGRPC        bool   `yaml:"use_grpc" json:"use_grpc"`
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
func (c *Config) GetGRPCServerHost() string {
	return c.GRPCServerHost
}
func (c *Config) IsUseGRPC() bool {
	return c.UseGRPC
}

func (c *Config) GetKey() string {
	return c.Key
}

func (c *Config) GetRateLimit() int {
	return c.RateLimit
}

func (c *Config) GetCryptoKey() string {
	return c.CryptoKey
}

func NewConfigCommand() (cf interfaces.ConfigAgent) {
	config := new(Config)

	// Чтение флагов
	flag.StringVar(&config.ServerHost, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.GRPCServerHost, "grpc-addr", "localhost:50051", "grpc address and port to run server")
	flag.IntVar(&config.ReportInterval, "r", reportInterval, "report interval")
	flag.IntVar(&config.PollInterval, "p", pollInterval, "poll interval")
	flag.StringVar(&config.Key, "k", "", "key")
	flag.IntVar(&config.RateLimit, "l", rateLimit, "rateLimit")
	flag.StringVar(&config.CryptoKey, "crypto-key", "", "crypto-key")
	flag.BoolVar(&config.UseGRPC, "use_grpc", false, "use_grpc")
	configFile := flag.String("c", "", "path to config file")

	flag.Parse()

	// Обработка переменных окружения
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

	if UseGRPC := os.Getenv("USE_GRPC"); UseGRPC != "" {
		UseGRPCBool, err := strconv.ParseBool(UseGRPC)
		if err == nil {
			config.UseGRPC = UseGRPCBool
		}
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		config.Key = envKey
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		envRateLimit, err := strconv.Atoi(envRateLimit)
		if err == nil {
			config.RateLimit = envRateLimit
		}
	}

	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		config.CryptoKey = envCryptoKey
	}

	// Если файл конфигурации указан, загружаем его
	if *configFile != "" {
		file, err := os.Open(*configFile)
		if err != nil {
			fmt.Println("Ошибка открытия файла конфигурации:", err)
			return nil
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Ошибка чтения файла конфигурации:", err)
			return nil
		}

		// Парсим JSON из файла и применяем
		var fileConfig Config
		err = json.Unmarshal(data, &fileConfig)
		if err != nil {
			fmt.Println("Ошибка парсинга JSON из конфигурационного файла:", err)
			return nil
		}

		// Перезаписываем только пустые поля, чтобы сохранить приоритет флагов и переменных окружения
		if config.ReportInterval == reportInterval {
			config.ReportInterval = fileConfig.ReportInterval
		}
		if config.PollInterval == pollInterval {
			config.PollInterval = fileConfig.PollInterval
		}
		if config.ServerHost == "localhost:8080" {
			config.ServerHost = fileConfig.ServerHost
		}
		if config.GRPCServerHost == "localhost:50051" {
			config.ServerHost = fileConfig.ServerHost
		}
		if config.Key == "" {
			config.Key = fileConfig.Key
		}
		if config.RateLimit == rateLimit {
			config.RateLimit = fileConfig.RateLimit
		}
		if config.CryptoKey == "" {
			config.CryptoKey = fileConfig.CryptoKey
		}
	}

	// Формируем ServerHost
	config.ServerHost = fmt.Sprintf("http://%s", config.ServerHost)
	return config
}
