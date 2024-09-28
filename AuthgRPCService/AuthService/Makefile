.PHONY: build up logs down

logs:
	docker-compose logs -f

build:
	go build -v ./cmd/AuthService

up:
	docker-compose up -d --build

down:
	docker-compose down

.DEFAULT_GOAL := up