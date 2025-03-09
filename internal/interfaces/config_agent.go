package interfaces

type ConfigAgent interface {
	GetReportInterval() int
	GetPollInterval() int
	GetServerHost() string
	GetGRPCServerHost() string
	GetKey() string
	GetRateLimit() int
	GetCryptoKey() string
}
