.ONESHELL:
TAB=echo "\t"
CURRENT_DIR = $(shell pwd)

help:
	@$(TAB) make up-server - запустить сервер
	@$(TAB) make up-agent - запустить агент
	@$(TAB) make doc - генерация документации
	@$(TAB) make cover-percent - процент покрытия тестами\(читаем из фаила отчета\)
	@$(TAB) make cover - отчет покрытия тестами
	@$(TAB) make cover-save - сохранить отчет покрытия тестами
	@$(TAB) make save-mem-prof - сохранить фаил профаилинга


up-server:
	go run ./cmd/server/main.go -a=localhost:41839 -d=postgresql://developer:developer@localhost:5432/postgres?sslmode=disable -k=secretkey

up-agent:
	go run ./cmd/agent/main.go -a=localhost:41839 -k=secretkey

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

doc:
	swag init -g ./cmd/server/main.go

cover:
	go test -cover ./...

cover-save:
	go test -coverprofile=coverage.out ./...

cover-percent:
	go tool cover -func=coverage.out | grep total

save-mem-prof:
	go test -bench=BenchmarkUpdateMemoryUsage -benchmem -memprofile profiles/base.pprof ./internal/services/agent
