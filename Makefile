.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	docker-compose up --remove-orphans app

down:
	docker-compose down

rebuild:
	docker-compose down && docker-compose build && docker-compose up

swag:
	@if ! [ -x "$$(command -v swag)" ]; then \
	  echo "swag is not installed or not found in PATH. Install it with 'go install github.com/swaggo/swag/cmd/swag@latest'"; \
	  exit 1; \
	fi
	swag init -g cmd/app/main.go
