package interfaces

import "time"

type ConfigServer interface {
	GetStoreInterval() time.Duration
	GetFileStoragePath() string
	IsRestore() bool
	GetFlagRunAddr() string
	GetDBDsn() string
	GetKey() string
}
