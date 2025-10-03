# Accounts Service

Manages user accounts, authentication and authorization for the SOA Social Network. Issues JWT access tokens and manages long-lived API tokens. Persists account data in PostgreSQL.

- Language: Go
- Storage: PostgreSQL (migrations under db/migrations)
- RPC: gRPC
- Auth: JWT (Ed25519), API tokens

## Responsibilities

- Register and manage user accounts and profiles
- Authenticate users and issue JWTs
- Create and validate long-lived API tokens
- Outbox pattern support for emission of domain events (e.g., registrations)

## gRPC API

Protocol buffers defined under services/accounts/proto/service.proto.
Generated code lives in the same directory.

Consumers: Gateway service and potentially other internal services.

## Environment variables

- ACCOUNTS_SERVICE_PORT: gRPC port to listen on
- DB_HOST: PostgreSQL host (e.g., accounts-postgres)
- DB_USER: PostgreSQL username
- DB_PASSWORD: PostgreSQL password
- DB_POOL_SIZE: Connection pool size (e.g., 5)
- JWT_ED25519_PRIVATE_KEY: Hex-encoded Ed25519 private key used to sign JWTs

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
export ACCOUNTS_SERVICE_PORT=50051
export DB_HOST=localhost
export DB_USER=...; export DB_PASSWORD=...
export DB_POOL_SIZE=5
export JWT_ED25519_PRIVATE_KEY=...

go run ./services/accounts/cmd
```

## Development

- Regenerate protobuf stubs: `make proto`
- Unit tests: `make test`

## Directory structure

- cmd/: service entrypoint
- internal/
  - models/: domain models (account, user, profile, api token, outbox)
  - repo/: repository interfaces and wiring
  - server/: gRPC server setup
  - service/: business logic, interceptors, config, outbox job
  - soajwtissuer/: JWT issuing helpers
  - storage/postgres/: PG implementation and tests
- db/migrations/: SQL schema migrations
- proto/: protobuf definitions and generated code

## Security notes

- Keep JWT private key secret and rotate regularly.
- Ensure DB credentials are provisioned securely.
