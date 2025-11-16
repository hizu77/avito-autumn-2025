# Быстрый старт (Docker)

**Требования:** Docker Desktop / Docker Engine (Compose v2), свободные порты `5432` (Postgres) и `8080` (HTTP API)

Запуск бд + миграций + приложения 

```bash
docker compose up -d
```

Сервис будет доступен по адресу:

```
http://localhost:8080
```