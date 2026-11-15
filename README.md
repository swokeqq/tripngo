# TripGo — trip-service

MVP сервиса поездок такси: создать поездку, получить её и завершить. HTTP API
по контракту OpenAPI, данные в PostgreSQL. 

## Что сделано в работе 1

- HTTP-сервер на `net/http` + `chi`: таймауты, graceful shutdown по `SIGINT` / `SIGTERM`.
- `/health` и `/ready`.
- Типы и серверный интерфейс сгенерированы из `contracts/openapi/trip-service.openapi.yaml`.
- Миграции `trips` и `trip_status_history` с `up` и `down`.
- Репозитории на `pgx` + `squirrel`, все запросы параметризованы.
- Контекстный менеджер транзакций: создание поездки атомарно пишет в две таблицы.
- Ручки `createTrip`, `getTrip`, `finishTrip`. Ошибки отдаются в `application/problem+json`.
- Запрет двух активных поездок на водителя и защита от двойного завершения, обе держатся под конкурентной нагрузкой.
- Задание со звёздочкой: многостадийный `Dockerfile`, итоговый образ на `distroless`, запуск не от root.

## Требования

- Go 1.25+
- [`tripgoctl`](https://github.com/course-go-autumn-2026/course-infra) — локальное окружение с PostgreSQL
- [`goose`](https://github.com/pressly/goose) в `PATH` — миграции:
  `go install github.com/pressly/goose/v3/cmd/goose@latest`
- `curl`, по желанию `psql` для проверок
- Docker — только для запуска в контейнере

## Запуск

```bash
tripgoctl cluster start        # один раз на машине
tripgoctl environment start    # PostgreSQL для этой работы, создаёт .env
make migrate                   # накатить миграции
make run                       # запустить сервис на HTTP_ADDR (по умолчанию :8080)
```

`Makefile` подхватывает `.env` сам.

| Цель | Что делает |
|---|---|
| `make run` | запустить сервис |
| `make migrate` | накатить миграции |
| `make migrate-down` | откатить последнюю миграцию |
| `make generate` | перегенерировать код по OpenAPI в `internal/generated` |

Тестов в этой работе нет, они появятся в работе 2.

## Docker

[`deploy/Dockerfile`](deploy/Dockerfile) собирает сервис в две стадии:

- **сборка** на `golang:1.27-alpine`: статический бинарь (`CGO_ENABLED=0`,
  `-trimpath -ldflags="-s -w"`). Зависимости скачиваются отдельным слоем, поэтому
  правка кода не сбрасывает их кеш;
- **итоговый образ** на `gcr.io/distroless/static-debian13:nonroot`: внутри только
  бинарь, без исходников, тулчейна, кеша сборки и шелла. Процесс запускается от
  пользователя `nonroot`.

[`.dockerignore`](.dockerignore) не пускает в контекст сборки `.git`, `.env`,
`.tripgo/`, `bin/`, тесты и файлы редактора.

**Размер итогового образа: 22.5 МБ.**

Сборка выполняется из корня репозитория:

```bash
docker build -f deploy/Dockerfile -t trip-service:lab01 .
```

Запуск рядом с PostgreSQL из `tripgoctl`. Бд доступна на `localhost` хоста,
поэтому контейнеру нужна сеть хоста. Конфигурация передаётся при запуске, в
образ `.env` не попадает:

```bash
docker run --rm --network host --env-file .env trip-service:lab01
```

Проверка: `curl localhost:8080/ready` → `200`. `docker stop` отправляет
`SIGTERM`, и сервис завершается штатно, через graceful shutdown.

## Переменные окружения

Обязательные параметры проверяются на старте: при ошибке сервис не запускается и
пишет причину в лог. Шаблон — [`.env.example`](.env.example).

| Переменная | По умолчанию | Описание |
|---|---|---|
| `HTTP_ADDR` | — (обязательная) | адрес HTTP-сервера, например `:8080` |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error` |
| `SHUTDOWN_TIMEOUT` | `10s` | сколько ждать завершения активных запросов при остановке |
| `HTTP_READ_HEADER_TIMEOUT` | `2s` | таймаут чтения заголовков запроса |
| `HTTP_READ_TIMEOUT` | `5s` | таймаут чтения всего запроса |
| `HTTP_WRITE_TIMEOUT` | `10s` | таймаут записи ответа |
| `HTTP_IDLE_TIMEOUT` | `30s` | сколько держать keep-alive соединение без запросов |
| `DATABASE_URL` | — (обязательная) | строка подключения к PostgreSQL |
| `DATABASE_MAX_CONNS` | `10` | максимум соединений в пуле |
| `DATABASE_MIN_CONNS` | `2` | минимум соединений в пуле |
| `DATABASE_MAX_CONN_LIFETIME` | `30m` | время жизни соединения |
| `DATABASE_CONNECT_TIMEOUT` | `5s` | таймаут проверки бд (`Ping`) на старте |
| `DATABASE_QUERY_TIMEOUT` | `3s` | таймаут одного запроса к бд |

## API

Источник — [`contracts/openapi/trip-service.openapi.yaml`](contracts/openapi/trip-service.openapi.yaml).

| Метод | Путь | Ответы |
|---|---|---|
| `POST` | `/api/v1/trips` | `201` + `Location`, `400`, `409 driver_busy`, `500` |
| `GET` | `/api/v1/trips/{tripId}` | `200`, `400`, `404 trip_not_found`, `500` |
| `POST` | `/api/v1/trips/{tripId}/finish` | `200`, `400`, `404 trip_not_found`, `409 trip_completed`, `500` |
| `GET` | `/health` | `200`, пока процесс жив; в бд не ходит |
| `GET` | `/ready` | `200`, если бд доступна, иначе `503` |

Пример:

```bash
curl -i -X POST localhost:8080/api/v1/trips -H 'content-type: application/json' \
  -d '{"user_id":"5cb72c04-7650-45c9-a79b-bcdba0631e0c",
       "driver_id":"8860b315-ec86-42eb-a17c-7c163d721ff5",
       "start_point":{"latitude":59.9398,"longitude":30.3146},
       "end_point":{"latitude":59.9290,"longitude":30.3626},
       "price":1450}'
```

## Структура

```
cmd/trip-service/           точка входа: конфиг, логгер, сигналы
internal/app/               сборка зависимостей, запуск и остановка сервера
internal/config/            конфигурация из env
internal/domain/            модели и доменные ошибки
internal/handler/           HTTP-ручки, валидация, problem+json
internal/repository/postgres/  пул, репозитории, менеджер транзакций
internal/generated/         код из OpenAPI (не редактируется руками)
migrations/                 миграции goose
deploy/                     Dockerfile
contracts/                  копия контрактов курса (не редактируется)
```

## Решения

### Менеджер транзакций

[`internal/repository/postgres/tx.go`](internal/repository/postgres/tx.go).
Интерфейс `Do(ctx, fn func(ctx) error) error` объявлен в пакете `handler`, где
он используется. Поэтому хендлеры не знают ни про `pgx`, ни про пул.

- `Do` открывает транзакцию, кладёт её в `context` под неэкспортируемым ключом
  и вызывает `fn` с новым контекстом.
- `fn` вернула `nil` → `COMMIT`. Вернула ошибку → `ROLLBACK`, и наружу уходит
  исходная ошибка: `errors.Is` в хендлере по-прежнему её распознаёт. Если упал
  и сам откат, обе ошибки склеиваются через `errors.Join`.
- Паника внутри `fn` → `ROLLBACK`, паника летит дальше.
- Откат выполняется с `context.WithoutCancel`: он доходит до бд, даже если
  клиент уже оборвал запрос.
- Вложенный `Do` находит транзакцию в контексте и переиспользует её. Commit или
  rollback делает только внешний `Do`.
- Репозитории берут исполнителя через `executorFromContext`: есть транзакция в
  контексте — работают через неё, нет — через пул. Транзакция не передаётся
  аргументом методов. Вызов репозитория вне `Do` выполняется в автокоммите.

### Уровень изоляции: Read Committed

Это уровень PostgreSQL по умолчанию. От гонок здесь защищают не уровень
изоляции, а ограничения и условия в запросах (см. ниже), поэтому повышать его не
нужно. `Serializable` тоже закрыл бы гонки, но при конфликте возвращает ошибку
`40001`, и тогда нужен цикл повторов, а это лишняя сложность для MVP.

### Запрет двух активных поездок у водителя

Частичный уникальный индекс в миграции:

```sql
CREATE UNIQUE INDEX trips_driver_active_unique_idx ON trips (driver_id) WHERE status = 'active';
```

Проверка `SELECT` перед `INSERT` не работает: два параллельных запроса оба
увидят, что водитель свободен. Индекс делает проверку и запись одной операцией
в бд. Второй параллельный `INSERT` ждёт на индексе и после коммита первого
получает `23505`. Репозиторий распознаёт `23505` **именно на этом индексе**
(по `ConstraintName`) и возвращает `domain.ErrDriverBusy`, то есть `409
driver_busy`. Коллизия первичного ключа даст `500`, а не «водитель занят».
Индекс частичный, поэтому завершённые поездки не мешают водителю взять новую.

Транзакция в создании нужна для атомарности двух вставок (`trips` и
`trip_status_history`), а не для защиты от гонки.

Проверено: 20 параллельных запросов на одного водителя дают ровно один `201`.

### Защита от двойного завершения

Условие на статус стоит прямо в `UPDATE`:

```sql
UPDATE trips SET status = 'completed', finished_at = $1, updated_at = $1
WHERE id = $2 AND status = 'active'
RETURNING ...
```

Проверка и запись — одна команда. Второй параллельный `UPDATE` ждёт блокировку
строки, после коммита первого перепроверяет `WHERE`, видит `completed` и
обновляет 0 строк. Тогда репозиторий дочитывает поездку, чтобы различить «нет
такой» (`404`) и «уже завершена» (`409 trip_completed`). `finished_at` не
перезаписывается.

Проверено: 20 параллельных `finish` одной поездки дают один `200` и
девятнадцать `409`, в журнале одна запись `active → completed`.

### Ошибки

Доменные ошибки превращаются в HTTP-ответы в одном месте —
`handler.writeError`. Все виды ошибок (`type`, `title`, `status`, `code`)
описаны таблицей в [`internal/handler/errors.go`](internal/handler/errors.go) и
совпадают с контрактом. Для `500` настоящая ошибка пишется в лог, а клиенту
уходит общий текст. Ошибки разбора параметров пути и заголовков (кривой
`tripId`) перехватывает `ParamErrorHandler`, поэтому они тоже отдаются в
problem+json, без внутренних имён Go-типов.

### Время

`started_at`, `finished_at` и `changed_at` в журнале выставляет сервис, один
раз на операцию, поэтому поездка и запись в журнале получают одно и то же
время. Наружу время всегда отдаётся в UTC.

## Ограничения

- Бизнес-логика пока живёт в хендлерах, без сервисного слоя. Слои — работа 2.
- Тестов нет — работа 2.
- Отсутствующее в теле поле координат неотличимо от `0`: без `start_point`
  поездка начнётся в точке `(0, 0)`. Отсутствующие `user_id` и `driver_id`
  ловятся (это нулевой UUID).
