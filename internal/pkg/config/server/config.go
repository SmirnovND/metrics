package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"io"
	"os"
	"strconv"
	"time"
)

const DefaultStoreInterval = 300

type Config struct {
	StoreInterval   int    `json:"store_interval"`
	FileStoragePath string `json:"store_file"`
	Restore         bool   `json:"restore"`
	FlagRunAddr     string `json:"address"`
	DBDsn           string `json:"database_dsn"`
	Key             string `json:"key"`
	CryptoKey       string `json:"crypto_key"`
}

func (c *Config) GetStoreInterval() time.Duration {
	return time.Second * time.Duration(c.StoreInterval)
}

func (c *Config) GetFileStoragePath() string {
	return c.FileStoragePath
}

func (c *Config) IsRestore() bool {
	return c.Restore
}

func (c *Config) GetFlagRunAddr() string {
	return c.FlagRunAddr
}

func (c *Config) GetDBDsn() string {
	return c.DBDsn
}

func (c *Config) GetKey() string {
	return c.Key
}

func (c *Config) GetCryptoKey() string {
	return c.CryptoKey
}

func NewConfigCommand() (cf interfaces.ConfigServerInterface) {
	config := new(Config)

	flag.IntVar(&config.StoreInterval, "i", DefaultStoreInterval, "")
	flag.StringVar(&config.DBDsn, "d", "", "db dsn")
	flag.StringVar(&config.FileStoragePath, "f", "./internal/resource/storage.json", "file storage path")
	flag.BoolVar(&config.Restore, "r", true, "Restore")
	flag.StringVar(&config.FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&config.Key, "k", "", "key")
	flag.StringVar(&config.CryptoKey, "crypto-key", "", "crypto-key")
	configFile := flag.String("c", "", "path to config file")

	flag.Parse()

	// Read config from environment variables
	if StoreInterval := os.Getenv("STORE_INTERVAL"); StoreInterval != "" {
		StoreIntervalInt, err := strconv.Atoi(StoreInterval)
		if err == nil {
			config.StoreInterval = StoreIntervalInt
		}
	}

	if FileStoragePath := os.Getenv("FILE_STORAGE_PATH"); FileStoragePath != "" {
		config.FileStoragePath = FileStoragePath
	}

	if Dsn := os.Getenv("DATABASE_DSN"); Dsn != "" {
		config.DBDsn = Dsn
	}

	if Restore := os.Getenv("RESTORE"); Restore != "" {
		RestoreBool, err := strconv.ParseBool(Restore)
		if err == nil {
			config.Restore = RestoreBool
		}
	}

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		config.FlagRunAddr = envRunAddr
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		config.Key = envKey
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

		// Check if the values are empty and only update empty fields
		if config.StoreInterval == 0 {
			config.StoreInterval = fileConfig.StoreInterval
		}
		if config.FileStoragePath == "" {
			config.FileStoragePath = fileConfig.FileStoragePath
		}
		if config.FlagRunAddr == "" {
			config.FlagRunAddr = fileConfig.FlagRunAddr
		}
		if config.DBDsn == "" {
			config.DBDsn = fileConfig.DBDsn
		}
		if config.Key == "" {
			config.Key = fileConfig.Key
		}
		if config.CryptoKey == "" {
			config.CryptoKey = fileConfig.CryptoKey
		}
	}

	return config
}
