# Repository Guidelines

## Project Structure & Module Organization

This Go API service uses Fiber, PostgreSQL, pgx, and SQLC. The entry point is `cmd/main.go`, which wires config, database access, middleware, services, repositories, and HTTP handlers.

- `internal/<domain>/`: feature modules such as `auth`, `user`, `center`, and `appointment`; each generally contains handlers, services, repositories, and middleware.
- `internal/domain/`: shared domain models.
- `internal/primitive/`: small shared primitive types, such as roles.
- `pkg/`: reusable helpers for errors, Fiber context values, and logging middleware.
- `db/queries/`: SQLC query definitions.
- `db/sqlc/`: generated SQLC Go code. Regenerate it; do not hand-edit generated files.
- `db/migrations/`: numbered PostgreSQL migration files.

## Build, Test, and Development Commands

- `make run`: run the API with `go run ./cmd/main.go`.
- `make docker-up`: start services with Docker Compose using `docker compose up -d --build`.
- `make migrate-up`: apply all pending database migrations.
- `make migrate-down`: roll back migrations.
- `make migrate-create name=create_example_table`: create a sequential SQL migration pair.
- `make sqlc-gen`: regenerate `db/sqlc` from `db/queries` and `db/migrations`.
- `go test ./...`: run the full Go test suite.

The Makefile loads `.env`, and required runtime variables are defined in `cmd/main.go` (`APP_PORT`, database settings, JWT settings, Mailjet settings, and `VERIFICATION_CODE_TTL`).

## Coding Style & Naming Conventions

Use standard Go formatting: run `gofmt` on changed Go files and keep imports organized with `go fmt ./...` or editor tooling. Package names should be short, lowercase, and domain-oriented. Keep feature code inside its owning `internal/<domain>` package, and preserve the existing handler-service-repository layering.

## Testing Guidelines

Place tests next to the code under test using Go’s `*_test.go` convention. Prefer table-driven tests for service and middleware logic. Repository tests should use an isolated PostgreSQL database or transaction cleanup. Run `go test ./...` before opening a PR.

## Commit & Pull Request Guidelines

Recent history uses short, imperative messages, often with prefixes such as `feature:` and `refactor:`. Keep commits focused, for example `feature: add appointment cancellation` or `refactor: move verification flow`.

Pull requests should include a concise summary, testing notes, migration notes when database files change, and example requests for API behavior when useful. If SQL queries or migrations change, include regenerated SQLC output in the same PR.

## Security & Configuration Tips

Do not commit real secrets from `.env`. Keep JWT secrets, database credentials, and Mailjet keys local or in deployment secret storage. When adding config, update `cmd/main.go` and document the new variable in PR notes.
