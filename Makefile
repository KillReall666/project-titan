.PHONY: up down rebuild logs

up:
	docker compose up --build

down:
	docker compose down

rebuild:
	docker compose build --no-cache

logs:
	docker compose logs -f
