.ONESHELL:
TAB=echo "\t"
CURRENT_DIR = $(shell pwd)

help:
	@$(TAB) make up-server - запустить сервер
	@$(TAB) make up-agent - запустить агент

up-server:
	go run ./cmd/server/main.go -a=localhost:41839 -d=postgresql://developer:developer@localhost:5432/postgres?sslmode=disable

up-agent:
	go run ./cmd/agent/main.go -a=localhost:41839

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)
