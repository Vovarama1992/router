.PHONY: refresh full-refresh build up down logs app-logs commit migrate db

APP_NAME=router
DB_NAME=router
DB_USER=router
DB_CONTAINER=router-db-1
APP_CONTAINER=router-router-1

# --- быстрый деплой приложения ---
refresh:
	git pull origin master
	docker compose build router
	docker compose stop router
	docker compose up -d --no-deps router
	docker compose logs -f router

# --- полный рефреш (без удаления volumes) ---
full-refresh:
	git pull origin master
	docker compose down
	docker compose build --no-cache
	docker compose up -d
	$(MAKE) migrate
	docker compose logs -f router

# --- сборка контейнеров ---
build:
	docker compose build

# --- поднять сервисы ---
up:
	docker compose up -d

# --- остановить сервисы ---
down:
	docker compose down

# --- логи всех сервисов ---
logs:
	docker compose logs -f

# --- логи только backend ---
app-logs:
	docker logs --tail=100 -f $(APP_CONTAINER)

# --- Git ---
commit:
	git add .
	git commit -m "$${m:-update}"
	git push origin master

# --- применить миграции ---
migrate:
	cat migrations/*.sql | docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

# --- зайти в PostgreSQL ---
db:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)