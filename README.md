# 🐌 SlugKiller — High-Performance URL Shortener

**SlugKiller** — это высокопроизводительный сервис сокращения ссылок на Go, поддерживающий хранение данных как в PostgreSQL, так и во встроенном In-Memory-хранилище, с кешированием через Redis и лаконичным REST API.

---

## 🚀 Возможности

- **Генерация уникальных коротких ссылок**
  - Длина ссылки — 10 символов (`a-zA-Z0-9_`).
  - Криптографически безопасная генерация (`crypto/rand`).
  - Надёжный механизм обработки коллизий.

- **Гибкие хранилища данных**
  - PostgreSQL для стабильного и надёжного хранения.
  - In-Memory storage для разработки и тестов.

- **Высокая производительность**
  - Кеширование ссылок с помощью Redis.

- **Чистая архитектура**
  - Чёткое разделение на слои (`Handler → Service → Repository`).
  - Dependency Injection через интерфейсы.

- **Расширенное логирование и мониторинг**
  - Подробное логирование через zerolog.

---

## 🛠️ Стек технологий

| Назначение       | Технология                 |
|------------------|----------------------------|
| Язык             | Go (1.21+)                 |
| HTTP-фреймворк   | Gin                        |
| База данных      | PostgreSQL, In-Memory      |
| Кеширование      | Redis                      |
| Логирование      | zerolog                    |
| Тестирование     | testify/mock               |
| Контейнеризация  | Docker, Docker Compose     |
| Миграции         | golang-migrate             |

---
## 📦 Структура проекта

```
SlugKiller/
├── cmd/
│   └── server/               # Точка входа (main.go)
│       └── main.go
├── config/                   # Конфигурация (чтение переменных окружения)
├── infrastructure/
│   ├── db/                   # Инициализация и подключение к PostgreSQL
│   └── redis/                # Инициализация и подключение к Redis
├── internal/
│   ├── cache/                # Работа с Redis (интерфейсы и реализация)
│   ├── handler/              # HTTP-обработчики (используется gin)
│   ├── model/                # Общие структуры данных
│   ├── repository/           # Интерфейсы репозиториев (URL, Slug и др.)
│   ├── server/               # Запуск и настройка HTTP-сервера
│   ├── service/              # Бизнес-логика
│   ├── storage/              # Реализации хранилищ
│   │   ├── mem/              # In-memory реализация
│   │   └── pg/               # PostgreSQL реализация
│   ├── init_storage.go       # Инициализация хранилища
│   └── init_storage_test.go  # Тесты инициализации хранилища
├── internal/tests/mocks/     # Моки для unit-тестов (генерируются mockery)
├── migrations/               # SQL-миграции базы данных
├── pkg/
│   └── logger/               # Обертка над zerolog
├── .env                      # Переменные окружения
├── docker-compose.yml        # Docker Compose конфигурация
├── Dockerfile                # Dockerfile для сервиса
├── Dockerfile.migrations     # Dockerfile для миграций
├── Dockerfile.tests          # Dockerfile для тестов
├── go.mod
├── go.sum
├── Makefile                  # Команды сборки, запуска и тестов
└── README.md                 # Документация проекта

```


---

## ⚙️ Установка и настройка

### Клонирование репозитория

```bash
git clone https://github.com/Thoustick/SlugKiller.git
cd SlugKiller
```

# HTTP Server
HTTP_ADDR=:8080

# Storage (postgres или memory)
STORAGE_TYPE=memory

# PostgreSQL
DB_HOST=postgres
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=slugkiller
DATABASE_URL=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable
DB_TIMEOUT_SECONDS=15

# Redis
REDIS_HOST=redis:6379
REDIS_PASSWORD=
REDIS_DB=0
CACHE_TTL_HOURS=0

# Slug settings
SLUG_LENGTH=10
MAX_ATTEMPTS=5

# Logger
LOG_LEVEL=debug|info|warn|error|fatal


## 🛠️ Запуск

### 🚧 Run with In-Memory Storage

```bash
docker-compose up --build
```

### 🧪 Run with PostgreSQL

```bash
STORAGE_TYPE=postgres docker-compose --profile postgres up --build
```
### Run Tests

```bash
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
```


## 📫 API Endpoints

### 1. POST `/shorten`

Сократить ссылку:

```http
POST /shorten
Content-Type: application/json

{
  "url": "https://example.com/very-long-url"
}
```

**Ответ:**

```json
{
  "slug": "AbC12_xYZ3"
}
```

### 2. GET `/{slug}`

Перенаправление по сокращённой ссылке:

```http
GET /AbC12_xYZ3
```

**Ответ:**

```
302 Found → Location: https://example.com
```

## ✅ Локальные Тесты

```bash
go test -v -cover ./...
```

Генерация отчёта покрытия:

```bash
go tool cover -html=coverage.out -o coverage.html
```

### 🛠️ Makefile Commands

| Команда         | Описание                                  |
|------------------|--------------------------------------------|
| `make build`     | Сборка всех Docker-контейнеров            |
| `make run`       | Запуск сервиса с In-Memory хранилищем     |
| `make run-pg`    | Запуск сервиса с PostgreSQL               |
| `make test`      | Запуск всех unit-тестов                   |
| `make migrate`   | Применение миграций к базе данных         |
| `make clean`     | Удаление всех контейнеров и образов       |
