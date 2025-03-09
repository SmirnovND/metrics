package interfaces

import "time"

type ConfigServerInterface interface {
	GetStoreInterval() time.Duration
	GetFileStoragePath() string
	IsRestore() bool
	GetFlagRunAddr() string
	GetDBDsn() string
	GetKey() string
	GetCryptoKey() string
	GetGRPCAddr() string
}
