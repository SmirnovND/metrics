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
	@$(TAB) make cover-func - покрытие по функциям
	@$(TAB) make save-mem-prof - сохранить фаил профаилинга
	@$(TAB) make staticlint - статический анализатор - запуск


up-server:
	go run ./cmd/server/main.go -a=localhost:41839 -d=postgresql://developer:developer@localhost:5432/postgres?sslmode=disable -k=secretkey -crypto-key=./.cert/private.pem

up-agent:
	go run ./cmd/agent/main.go -a=localhost:41839 -k=secretkey -crypto-key=./.cert/public.pem

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

doc:
	swag init -g ./cmd/server/main.go

cover:
	go test -cover ./...

cover-save:
	go test -coverprofile=coverage.out ./...

cover-func:
	go tool cover -func=coverage.out

cover-percent:
	go tool cover -func=coverage.out | grep total

save-mem-prof:
	go test -bench=BenchmarkUpdateMemoryUsage -benchmem -memprofile profiles/base.pprof ./internal/services/agent

staticlint:
	go run cmd/staticlint/main.go ./...

