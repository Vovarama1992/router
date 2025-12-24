.PHONY: refresh build run restart stop status logs commit migrate db

APP_NAME=router
APP_BIN=./bin/router
ENV_FILE=.env

# --- быстрый деплой приложения ---
refresh:
	git pull origin master
	$(MAKE) build
	sudo systemctl restart $(APP_NAME)
	$(MAKE) logs

# --- сборка бинаря ---
build:
	mkdir -p bin
	go build -o $(APP_BIN) ./cmd

# --- запуск вручную (для отладки) ---
run:
	$(APP_BIN)

# --- systemd ---
restart:
	sudo systemctl restart $(APP_NAME)

stop:
	sudo systemctl stop $(APP_NAME)

status:
	sudo systemctl status $(APP_NAME) --no-pager

logs:
	sudo journalctl -u $(APP_NAME) -n 200 -f

# --- Git ---
commit:
	git add .
	git commit -m "$${m:-update}"
	git push origin master

# --- применить миграции ---
migrate:
	@set -a; . $(ENV_FILE); set +a; cat migrations/*.sql | psql "$$DATABASE_URL"

# --- зайти в PostgreSQL ---
db:
	@set -a; . $(ENV_FILE); set +a; psql "$$DATABASE_URL"