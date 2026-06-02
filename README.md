# go-auth

gRPC-сервис аутентификации: регистрация пользователей, вход и выдача JWT-токенов. Хранилище — SQLite, миграции — golang-migrate.

## Стек

Go · gRPC · SQLite (`mattn/go-sqlite3`) · golang-migrate · Taskfile · Docker

## API

| Метод | Запрос | Ответ |
|-------|--------|-------|
| `Register` | `email`, `password` | `user_id` |
| `Login` | `email`, `password`, `app_id` | `token` (JWT) |

## Конфигурация

Параметры задаются YAML-файлом (`config/`):

```yaml
env: "local"
storage_path: "./storage/auth.db"
token_ttl: 1h
grpc:
  port: 40049
  timeout: 10h
```

Путь к конфигу передаётся переменной окружения `CONFIG_PATH`.

## Запуск в Docker

```bash
cd docker
docker compose up --build
```

Сервис стартует на порту `40049`, миграции применяются автоматически при запуске контейнера. Данные сохраняются в Docker-томе `auth-db`.

```bash
docker compose down        # остановить
docker compose down -v     # остановить и удалить данные
```

## Локальный запуск

Требуется [Task](https://taskfile.dev).

```bash
task migrate   # применить миграции
task run       # собрать и запустить сервис
```

## Тесты

Интеграционные тесты выполняются против запущенного сервиса:

```bash
task test_migrate          # подготовить тестовые данные
task run                   # запустить сервис (в отдельном терминале)
go test ./tests/ -v
```

## Структура проекта

```
cmd/auth        точка входа сервиса
cmd/migrator    утилита применения миграций
internal/       gRPC-слой, бизнес-логика, хранилище, конфигурация
migrations/     SQL-миграции
tests/          интеграционные тесты
docker/         Dockerfile и docker-compose
```
