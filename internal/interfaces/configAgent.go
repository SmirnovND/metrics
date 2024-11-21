package interfaces

type ConfigAgent interface {
	GetReportInterval() int
	GetPollInterval() int
	GetServerHost() string
	GetKey() string
}
