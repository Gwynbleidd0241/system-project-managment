.PHONY: up logs down

logs:
	docker-compose logs -f

up:
	docker-compose up -d --build

down:
	docker-compose down

.DEFAULT_GOAL := up