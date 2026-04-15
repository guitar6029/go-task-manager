.PHONY: lint tidy build run test docker-build docker-up help

help:
	@echo "Available commands:"	
	@echo "	make lint			Run linter"
	@echo "	make tidy			Clean go.mod"
	@echo "	make build			Build API binary"
	@echo "	make run			Run API locally"

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./... -v
