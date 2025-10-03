# Posts Service

Provides profile pages, posts, and comments functionality for the SOA Social Network. Uses PostgreSQL for storage and exposes a gRPC API consumed by the Gateway.

- Language: Go
- Storage: PostgreSQL (migrations under db/migrations)
- RPC: gRPC
- Auth: JWT verification for protected operations (uses public key)

## Responsibilities

- Manage profile pages settings
- CRUD for posts
- CRUD for comments under posts
- Outbox for events sent to Stats (views, likes, comments, new posts)

## gRPC API

Protocol buffers defined under services/posts/proto/posts_service.proto.
Generated code lives in the same directory.

## Environment variables

- POSTS_SERVICE_PORT: gRPC port to listen on
- DB_HOST: PostgreSQL host (e.g., posts-postgres)
- DB_USER: PostgreSQL username
- DB_PASSWORD: PostgreSQL password
- DB_POOL_SIZE: Connection pool size (e.g., 5)
- JWT_ED25519_PUBLIC_KEY: Hex-encoded Ed25519 public key used to verify JWTs

## Database

- Docker Compose mounts db/migrations into the Postgres container for initialization.

## Run

Run with the full stack (recommended):

```bash
cp ../../test.env ../../.env   # from repo root
make -C ../../ run
```

Local run (requires PostgreSQL):

```bash
export POSTS_SERVICE_PORT=50052
export DB_HOST=localhost
export DB_USER=...; export DB_PASSWORD=...
export DB_POOL_SIZE=5
export JWT_ED25519_PUBLIC_KEY=...

go run ./services/posts/cmd
```

## Development

- Regenerate protobuf stubs: `make proto`
- Unit tests for Posts: `make test`

## Directory structure

- cmd/: service entrypoint
- internal/
  - models/: domain models (posts, comments, pages, outbox)
  - repo/: repositories and pagination helpers
  - server/: gRPC server setup
  - service/: business logic and jobs
  - storage/postgres/: PG implementation
- db/migrations/: SQL schema migrations
- proto/: protobuf definitions and generated code

## Notes

- Certain endpoints require JWT authentication. Ensure the public key matches the Accounts service private key counterpart.
