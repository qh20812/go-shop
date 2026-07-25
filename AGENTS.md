# go-shop

Two-package repo: `client/` (Next.js 16) + `server/` (Go 1.26, Gin).

## Client (`client/`)

- **Next.js 16.2.10** — read `client/AGENTS.md` first; breaking changes vs training data.
- **CSS**: plain CSS (no Tailwind). Styles in `app/globals.css`.
- **Proxy**: `next.config.ts` rewrites `/api/*` → `localhost:8080` (server Gin).
- **Commands**: `npm run dev` (port 3000), `npm run build`, `npm run start`, `npm run lint`. No typecheck or test scripts exist.

## Server (`server/`)

- **Go 1.26.3**, module `github.com/qh20812/go_shop_server`.
- **Framework**: Gin (`gin-gonic/gin` v1.12).
- **Database**: PostgreSQL (`pgx/v5` in go.mod). MongoDB driver (`mongo-driver/v2`) is in go.mod but unused — schema is PostgreSQL-only.
- **Schema**: `schema.sql` defines the `products` table (serial PK, name, price, description, image, created_at TIMESTAMPTZ). Load it after starting PostgreSQL.
- **Env**: `godotenv` loads `.env` automatically. Contains `PORT=8080` and `DATABASE_URL`.
- **DATABASE_URL format**: `postgres://user:password@host:5432/dbname?sslmode=disable`. `sslmode=disable` for local Docker (127.0.0.1). Remove for cloud DBs (Supabase, Neon).
- **Local PostgreSQL via Docker**:
  ```
  docker run -d --name pg-goshop -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=newpassword -e POSTGRES_DB=goshop -p 5432:5432 postgres:16
  docker exec -i pg-goshop psql -U postgres -d goshop < server/schema.sql
  ```
- **Commands**: `go run .` from `server/`. No test files exist.
- No Docker config, no CI, no `opencode.json`.
