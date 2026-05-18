# Go Beresin — Agent Guide

## Entrypoints

| Purpose | Command |
|---|---|
| API server | `cmd/api/main.go` |
| DB migrations | `cmd/migrate/main.go up` / `down` |
| Seed data | `cmd/seed/main.go` |

## Architecture (Layered)

```
cmd/api/main.go  →  internal/transport/router/router.go
                        ├── internal/transport/handler/   (HTTP only: parse → call service → respond)
                        ├── internal/transport/middleware/ (JWT auth, rate limiter)
                        ├── internal/service/             (business logic)
                        ├── internal/domain/              (DTOs, errors, UUID helpers)
                        └── internal/repository/           (sqlc-generated, DO NOT edit manually)
```

## Key Commands

```bash
go run cmd/api/main.go              # start server (port 8080, requires .env)
go run cmd/migrate/main.go up       # apply DB migrations (PostgreSQL + PostGIS)
go run cmd/seed/main.go             # seed sample data
go build ./...                      # verify compilation
go vet ./...                        # lint
go test ./internal/...              # integration tests (needs live Postgres + Redis)
```

## Testing

- Tests are **integration tests** — require real PostgreSQL and Redis running locally.
- Reads `.env` from project root. Redis defaults: `localhost:6379`, password `!Abcd1234`.
- Test setup in `auth_handler_test.go` calls `database.InitPool()` and connects to Redis directly.
- Test cleanup deletes specific Redis keys; does NOT clean DB records (manual cleanup via raw SQL in test).

## Codegen

- **sqlc**: source of truth is `internal/repository/query.sql`. Edit that, then run:
  ```bash
  $HOME/go/bin/sqlc generate
  ```
  Never edit `query.sql.go` or `models.go` manually.
- **Swagger**: regenerate after handler annotation changes:
  ```bash
  $HOME/go/bin/swag init -g cmd/api/main.go
  ```

## Conventions

- **Layered auth flow**: handler parses HTTP → calls `service.AuthService` interface → maps domain errors to HTTP status. No business logic in handlers.
- **UUID v7**: all new records use `domain.NewUUIDV7()` (time-ordered, better index perf). `google/uuid` v1.6.0.
- **INSERTs**: all sqlc INSERT queries pass `id` explicitly via named params (`sqlc.arg(id)`).
- **Response format**: `{"status":"success|error", "message":"...", "data":{...}}` — `Response` struct in `dummy_handler.go`.
- **JWT tokens**: access (15m) + refresh (30d). Refresh stored in Redis key `user:{id}:refresh`.
- **Rate limiting**: global 60 req/min, auth endpoints 5 req/min. IP blocked 15min after 5 failed logins.
- **Error mapping in handler**: `domain.ErrEmailConflict`/`ErrPhoneConflict` → 409, `ErrInvalidCredentials` → 401, `ErrIPBlocked` → 429, else 500.

## Secrets

- `.env` file (gitignored) or Doppler CLI. See `.env.example` for template.
- Test env vars from `../../../.env` relative to test file.
