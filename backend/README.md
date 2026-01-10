# ERP Backend (Go + Gin) Starter

Clean Architecture scaffold for ERP with Gin, PostgreSQL, sqlc, Viper, Zap, JWT auth, and golang-migrate migrations.

## Tech Choices
- Gin (fast, minimal, great middleware ecosystem) for HTTP.
- Zap for structured logging (lower allocs than Logrus).
- Viper for env-based config.
- pgx + sqlc (to be generated) for DB access.
- JWT (access/refresh) for auth; manual DI wiring.

## Layout
- `cmd/server` — entrypoint.
- `internal/config` — Viper config loader.
- `internal/logger` — Zap logger.
- `internal/pkg/database` — PostgreSQL connection helper.
- `internal/domain` — entities + repository interfaces (user, role, permission, product, inventory, audit, auth models).
- `internal/usecase` — business logic/services per feature.
- `internal/infrastructure/repository/postgres` — Postgres implementations (stubbed, wire sqlc later).
- `internal/transport/http` — Gin router, handlers, middleware.
- `migrations/` — golang-migrate files (add `.sql`).
- `docs/` — Swagger (swaggo) output.
- `internal/pkg/security` — bcrypt password hashing + JWT helpers.

## Running (dev)
```bash
cd backend
go mod tidy
# generate sqlc (optional, requires sqlc installed)
# sqlc generate

export DATABASE_URL="postgres://user:pass@localhost:5432/erp?sslmode=disable"
export JWT_ACCESS_SECRET="changeme-access"
export JWT_REFRESH_SECRET="changeme-refresh"
export CORS_ORIGINS="http://localhost:5173"
export COOKIE_SECURE="false"

# run migrations (example)
# migrate -path migrations -database "$DATABASE_URL" up

HTTP_PORT=8080 go run cmd/server/main.go
```

## Configuration
Environment variables (all optional):
- `APP_NAME` (default: `erp-backend`)
- `APP_ENV` (default: `development`)
- `HTTP_PORT` (default: `8080`)
- `HTTP_READ_TIMEOUT` (default: `5s`)
- `HTTP_WRITE_TIMEOUT` (default: `10s`)
- `SHUTDOWN_TIMEOUT` (default: `10s`)
- `DATABASE_URL` (required for DB connection)
- `JWT_ACCESS_SECRET` (required for auth)
- `JWT_REFRESH_SECRET` (required for auth)
- `JWT_ACCESS_TTL` (default: `15m`)
- `JWT_REFRESH_TTL` (default: `168h`)
- `CORS_ORIGINS` (default: `http://localhost:5173`)
- `COOKIE_SECURE` (default: `false`)

## API Sketch
- `POST /api/v1/auth/register|login|refresh`
- `GET /api/v1/health`
- Auth required (`Authorization: Bearer <token>`):
  - `POST /api/v1/users`, `GET /api/v1/users/:id`
  - `POST /api/v1/roles`, `GET /api/v1/roles`
  - `POST /api/v1/permissions`, `GET /api/v1/permissions`
  - `POST /api/v1/products`, `GET /api/v1/products`
  - `POST /api/v1/inventory`, `PATCH /api/v1/inventory/:id/adjust`, `GET /api/v1/inventory/product/:product_id`
  - `GET /api/v1/audit/logs?actor_id=:actor`
  - Swagger UI: `GET /docs/index.html` (uses `docs/swagger.json`)

## Next Steps
1. Replace manual sqlc stubs with generated code (`sqlc generate`) and hook up to real DB tests.
2. Add integration tests for repositories against a test database + migrations; extend unit tests coverage across use cases.
3. Add Swagger annotations in handlers and regenerate `docs/swagger.json` with swaggo CLI.
4. Implement refresh-token persistence/blacklist beyond in-memory and add rate limiting for auth endpoints.
