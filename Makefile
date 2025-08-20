ifneq (,$(wildcard .env))
include .env
export $(shell sed -n 's/^\([A-Za-z_][A-Za-z0-9_]*\)=.*/\1/p' .env)
endif

CONFIG_PATH ?= conf

DB_HOST := $(shell yq '.postgres.host' $(CONFIG_PATH)/config.yaml 2>/dev/null)
DB_PORT := $(shell yq '.postgres.port' $(CONFIG_PATH)/config.yaml 2>/dev/null)
DB_NAME := $(shell yq '.postgres.db'   $(CONFIG_PATH)/config.yaml 2>/dev/null)

DB_USER     ?= $(POSTGRES_USER)
DB_PASSWORD ?= $(POSTGRES_PASSWORD)

DB_URL       := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
DB_URL_NO_DB := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/postgres?sslmode=disable

APP := go run ./cmd/app
APP_ENV := CONFIG_PATH=$(CONFIG_PATH) POSTGRES_USER=$(DB_USER) POSTGRES_PASSWORD=$(DB_PASSWORD)

.PHONY: help
help: 
	@awk 'BEGIN {FS = ":.*##"; printf "\nДоступные команды:\n\n"} /^[a-zA-Z0-9_\-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 } END { print "" }' $(MAKEFILE_LIST)

.PHONY: deps
deps: 
	@command -v yq >/dev/null 2>&1    || (echo "📦 Устанавливаю yq..." && brew install yq)
	@command -v psql >/dev/null 2>&1  || (echo "📦 Устанавливаю libpq (psql)..." && brew install libpq && brew link --force libpq)

.PHONY: go-deps
go-deps:
	go mod tidy
	go mod download

.PHONY: createdb
createdb: deps 
	@[ -n "$(DB_HOST)" ] && [ -n "$(DB_PORT)" ] && [ -n "$(DB_NAME)" ] || (echo "❌ Проверь $(CONFIG_PATH)/config.yaml (postgres.host/port/db)"; exit 1)
	@[ -n "$(DB_USER)" ] && [ -n "$(DB_PASSWORD)" ] || (echo "❌ Укажи POSTGRES_USER/POSTGRES_PASSWORD в .env или окружении"; exit 1)
	@echo "🛠️  Проверяю наличие БД '$(DB_NAME)'..."
	@psql "$(DB_URL_NO_DB)" -v ON_ERROR_STOP=1 -tAc "SELECT 1 FROM pg_database WHERE datname='$(DB_NAME)'" | grep -q 1 \
		|| psql "$(DB_URL_NO_DB)" -v ON_ERROR_STOP=1 -c "CREATE DATABASE $(DB_NAME)"
	@echo "✅ База данных готова"

.PHONY: migrate-up
migrate-up: createdb 
	@$(APP_ENV) $(APP) migrate up

.PHONY: migrate-down
migrate-down:
	@$(APP_ENV) $(APP) migrate down

.PHONY: run
run: 
	@$(APP_ENV) $(APP) serve

.PHONY: up
up: deps go-deps migrate-up run
	@echo "✅ Всё готово!"

.PHONY: dsn
dsn: 
	@echo "$(DB_URL)"