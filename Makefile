.PHONY: run build test lint tidy docker-up docker-down

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test -race -cover ./...

lint:
	go vet ./...
	gofmt -l .

tidy:
	go mod tidy

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v
