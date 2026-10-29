# Users Service

Микросервис пользователей на Go (gRPC + PostgreSQL). Используется сервисом [tasks-service](https://github.com/AntonRadchenko/tasks-service): тот обращается к методу `GetUser`, чтобы проверить, что пользователь существует.

**Быстрый запуск и проверка всей связки из двух сервисов описаны в [README tasks-service](https://github.com/AntonRadchenko/tasks-service#readme).** Ниже то, что нужно для запуска только users-service.

**Стек:** Go 1.24, gRPC, GORM, PostgreSQL, bcrypt. Proto-контракты лежат в отдельном репозитории [project-protos](https://github.com/AntonRadchenko/project-protos) (v1.0.4) и подтягиваются автоматически через Go modules.

## Что нужно установить

- Go 1.24+
- Docker
- [grpcurl](https://github.com/fullstorydev/grpcurl) для проверки ручек (`brew install grpcurl`)

## Запуск базы данных (PostgreSQL)

```bash
docker run --name services-postgres \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=yourpassword \
    -e POSTGRES_DB=postgres \
    -p 5432:5432 \
    -d postgres:16
```

Параметры подключения (зашиты в `internal/database/db.go`):
- Host: localhost
- Port: 5432
- User: postgres
- Password: yourpassword
- DB: postgres
- SSL Mode: disable

## Создание таблицы

Миграций в проекте нет, таблица создаётся один раз вручную (подождите пару секунд после запуска контейнера):

```bash
docker exec -i services-postgres psql -U postgres -d postgres <<'SQL'
CREATE TABLE IF NOT EXISTS user_structs (
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT NOT NULL UNIQUE,
    password   TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
SQL
```

## Запуск сервиса

gRPC-сервер стартует на порту `50051`:

```bash
git clone https://github.com/AntonRadchenko/users-service.git
cd users-service
go run ./cmd/server
```

## gRPC методы

| Метод | Что делает |
|---|---|
| `CreateUser` | Создать пользователя (пароль хранится в виде bcrypt-хэша) |
| `GetUser` | Получить пользователя по id |
| `ListUsers` | Список всех пользователей |
| `UpdateUser` | Изменить `email` и/или `password` |
| `DeleteUser` | Удалить пользователя |

Валидация: пустые `email`/`password` и обновление без полей возвращают ошибку; повторный `email` даёт ошибку `email already exists`. Пароль в ответах не возвращается.

## Примеры запросов (grpcurl)

Для grpcurl нужны `.proto` файлы:

```bash
git clone --branch v1.0.4 --depth 1 https://github.com/AntonRadchenko/project-protos.git
cd project-protos
```

**1) Создать пользователя:**
```bash
grpcurl -plaintext -import-path . -proto proto/user/user.proto \
    -d '{"email": "user@example.com", "password": "123"}' \
    localhost:50051 user.UserService/CreateUser
```

**2) Получить пользователя по id:**
```bash
grpcurl -plaintext -import-path . -proto proto/user/user.proto \
    -d '{"id": 1}' \
    localhost:50051 user.UserService/GetUser
```

**3) Список всех пользователей:**
```bash
grpcurl -plaintext -import-path . -proto proto/user/user.proto \
    -d '{}' \
    localhost:50051 user.UserService/ListUsers
```

**4) Обновить пользователя:**
```bash
grpcurl -plaintext -import-path . -proto proto/user/user.proto \
    -d '{"id": 1, "email": "new@example.com"}' \
    localhost:50051 user.UserService/UpdateUser
```

**5) Удалить пользователя:**
```bash
grpcurl -plaintext -import-path . -proto proto/user/user.proto \
    -d '{"id": 1}' \
    localhost:50051 user.UserService/DeleteUser
```

## Структура проекта

```
cmd/server/main.go            # точка входа: БД → сервис → gRPC-сервер
internal/
├── database/db.go            # подключение к PostgreSQL (GORM)
├── user/
│   ├── model.go              # модель таблицы user_structs
│   ├── repository.go         # работа с БД
│   └── service.go            # бизнес-логика, валидация, хэширование пароля
└── transport/grpc/
    ├── server.go             # запуск gRPC-сервера на :50051
    └── handler.go            # реализация методов UserService
```

Слои: `handler` (gRPC) → `service` (логика) → `repository` (GORM).

## Остановка

В терминале с сервисом нажать `Ctrl+C`, затем:

```bash
docker stop services-postgres && docker rm services-postgres
```
