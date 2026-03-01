SHELL := /bin/bash

.PHONY: refresh build run restart stop status logs commit migrate db build-front

APP_NAME=router
BIN=./bin/router
GO=/usr/bin/go
ENV_FILE=.env

refresh:
	git pull origin master
	$(MAKE) build
	sudo systemctl restart $(APP_NAME)
	$(MAKE) logs

build:
	mkdir -p bin
	$(GO) build -o $(BIN) ./cmd

run:
	$(BIN)

restart:
	sudo systemctl restart $(APP_NAME)

stop:
	sudo systemctl stop $(APP_NAME)

status:
	sudo systemctl status $(APP_NAME) --no-pager

logs:
	sudo journalctl -u $(APP_NAME) -n 200 -f

commit:
	git add .
	git commit -m "$${m:-update}"
	git push origin master

migrate:
	@set -a; source $(ENV_FILE); set +a; \
	psql "$$DATABASE_URL" < migrations/001_create_peers.sql

db:
	@set -a; source $(ENV_FILE); set +a; \
	psql "$$DATABASE_URL"

build-front:
	cd ../router-front && \
	npm install && \
	npm run build

	rm -rf front-dist/*
	mkdir -p front-dist
	cp -r ../router-front/dist/* front-dist/