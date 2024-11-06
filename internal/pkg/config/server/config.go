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

func NewConfigCommand() (cf interfaces.ConfigServer) {
	config := new(Config)

	flag.IntVar(&config.StoreInterval, "i", DefaultStoreInterval, "")
	flag.StringVar(&config.DBDsn, "d", "", "db dsn")
	flag.StringVar(&config.FileStoragePath, "f", "./internal/resource/storage.json", "file storage path")
	flag.BoolVar(&config.Restore, "r", true, "Restore")
	flag.StringVar(&config.FlagRunAddr, "a", "localhost:8080", "address and port to run server")

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

	if FileStoragePath := os.Getenv("DATABASE_DSN"); FileStoragePath != "" {
		config.FileStoragePath = FileStoragePath
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

	return config
}
