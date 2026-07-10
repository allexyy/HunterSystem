# Hunter System

Task tracker в духе Solo Leveling: квесты, XP, статы, ранги охотника.

## Стек

- Go 1.24+
- Gin
- PostgreSQL
- sqlc (pgx/v5)
- Docker Compose
- Telegram Bot API

## Быстрый старт

1. Скопируй `.env.example` в `.env` и при желании поправь значения:

   ```bash
   cp .env.example .env
   ```

2. Установи `golang-migrate` (нужен для миграций):

   ```bash
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

   Убедись, что `$GOPATH/bin` (обычно `~/go/bin`) добавлен в `PATH`.

3. Подними Postgres:

   ```bash
   make up
   ```

4. Прогони миграции:

   ```bash
   make migrate-up
   ```

5. Проверь версию миграции (должна быть `1`, `dirty: false`):

   ```bash
   make migrate-version
   ```

6. Подтяни Go-зависимости:

   ```bash
   make tidy
   ```

7. Запусти API:

   ```bash
   make run-api
   ```

   Проверка: `curl localhost:8080/health` → `{"status":"ok"}`

## Структура

```
cmd/
    api/        — HTTP API (Gin)
    bot/        — Telegram bot
internal/
    habit/      — привычки (шаблоны для daily quest)
    quest/      — квесты (daily/main/boss)
    user/       — пользователи, статы, ранг
    reward/     — shop, inventory, achievements
    stats/      — начисление и агрегация характеристик
    telegram/   — хендлеры Telegram Bot API
pkg/
    database/   — подключение к Postgres (pgx pool)
migrations/     — SQL миграции (golang-migrate)
queries/        — SQL запросы для sqlc
```

## Полезные команды

Смотри `Makefile` — там есть `make up/down/logs`, `make migrate-*`,
`make sqlc-generate`, `make run-api/run-bot`.

## Roadmap

См. обсуждение в проекте — этапы 0–11, от core-квестов до Random Events и фронтенда.
