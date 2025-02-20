package server

import (
	"flag"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"os"
	"strconv"
	"time"
)

const DefaultStoreInterval = 300

type Config struct {
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	FlagRunAddr     string
	DBDsn           string
	Key             string
	CryptoKey       string
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

	flag.Parse()

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
		Restore, err := strconv.ParseBool(Restore)
		if err == nil {
			config.Restore = Restore
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

	return config
}
