.ONESHELL:
TAB=echo "\t"
CURRENT_DIR = $(shell pwd)

help:
	@$(TAB) make up-server - запустить сервер
	@$(TAB) make up-agent - запустить агент
	@$(TAB) make doc - генерация документации

up-server:
	go run ./cmd/server/main.go -a=localhost:41839 -d=postgresql://developer:developer@localhost:5432/postgres?sslmode=disable -k=secretkey

up-agent:
	go run ./cmd/agent/main.go -a=localhost:41839 -k=secretkey

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

doc:
	swag init -g ./cmd/server/main.go
