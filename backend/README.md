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
# install deps
# go mod tidy

export DATABASE_URL="postgres://user:pass@localhost:5432/erp?sslmode=disable"
export JWT_ACCESS_SECRET="changeme-access"
export JWT_REFRESH_SECRET="changeme-refresh"

HTTP_PORT=8080 go run cmd/server/main.go
```

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

## Next Steps
1. Generate SQLC code (schema + queries) into `internal/infrastructure/repository/postgres/sqlc` and implement repository methods.
2. Enforce RBAC permission checks in middleware and on routes; populate permissions per role.
3. Write migrations with golang-migrate for users/roles/permissions/products/inventory/audit tables (UUID PKs, audit fields, soft delete).
4. Add Swagger docs via swaggo annotations under handlers and generate to `docs/`.
5. Add unit tests for use cases and integration tests for repositories (70%+ coverage target).
