.ONESHELL:
TAB=echo "\t"
CURRENT_DIR = $(shell pwd)

help:
	@$(TAB) make up-server - запустить сервер
	@$(TAB) make up-agent - запустить агент

up-server:
	go run ./cmd/server/main.go -a=localhost:41839

up-agent:
	go run ./cmd/agent/main.go -a=localhost:41839

