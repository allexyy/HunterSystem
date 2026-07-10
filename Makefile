include .env
export

# Требуется локально установленный golang-migrate:
#   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# (бинарник должен попасть в $GOPATH/bin, добавленный в PATH)
MIGRATE := migrate -path=migrations -database "$(DATABASE_URL)"

.PHONY: up down restart logs \
	migrate-up migrate-down migrate-down-all migrate-force migrate-new migrate-version \
	sqlc-generate \
	run-api run-bot \
	tidy

## --- Docker / инфраструктура ---

up: ## Поднять Postgres в фоне
	docker compose up -d

down: ## Остановить и удалить контейнеры (данные в volume сохраняются)
	docker compose down

restart: down up ## Перезапустить окружение

logs: ## Логи Postgres
	docker compose logs -f postgres

## --- Миграции (требуется golang-migrate, см. комментарий выше) ---

migrate-up: ## Накатить все миграции
	$(MIGRATE) up

migrate-down: ## Откатить одну последнюю миграцию
	$(MIGRATE) down 1

migrate-down-all: ## Откатить вообще все миграции
	$(MIGRATE) down -all

migrate-force: ## Сбросить dirty-состояние. Использование: make migrate-force V=1
	$(MIGRATE) force $(V)

migrate-version: ## Показать текущую версию миграции
	$(MIGRATE) version

migrate-new: ## Создать новую пару миграций. Использование: make migrate-new NAME=add_shop
	migrate create -ext sql -dir migrations -seq $(NAME)

## --- sqlc ---

sqlc-generate: ## Сгенерировать Go-код из SQL запросов
	docker run --rm -v $(shell pwd):/src -w /src sqlc/sqlc generate

## --- Запуск сервисов локально ---

run-api:
	go run ./cmd/api

run-bot:
	go run ./cmd/bot

## --- Прочее ---

tidy:
	go mod tidy

lint:
	# Требуется локально установленный golangci-lint:
	golangci-lint run