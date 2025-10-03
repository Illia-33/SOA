# SOA Social Network

Service-oriented social network implemented in Go. The system is split into independent services and orchestrated with Docker Compose. A lightweight HTTP Gateway exposes a REST API to UI/clients, while internal services communicate via gRPC and Kafka.

- Language: Go
- Orchestration: Docker Compose
- Data stores: PostgreSQL (Accounts, Posts), ClickHouse (Stats)
- Messaging: Kafka (events for stats aggregation)
- HTTP: Gin (Gateway)
- RPC: gRPC (internal service-to-service)

## Services

- Gateway: Public REST API. Delegates requests to internal services, verifies JWT.
- Accounts: User accounts, authentication, authorization, long-lived API tokens.
- Posts: Profiles, pages, posts, comments.
- Stats: Views, likes, comments, registrations, and post analytics (ClickHouse + Kafka consumers).

See per-service READMEs for brief descriptions:
- services/gateway/README.md
- services/accounts/README.md
- services/posts/README.md
- services/stats/README.md

Diagrams are available under docs/ (C4/PlantUML sources).

## Repository layout

- deploy/: Docker and Compose definitions
- scripts/: Helper scripts invoked by Make targets
- services/
  - gateway/: HTTP entrypoint, REST to gRPC translation
  - accounts/: Accounts + auth service (PostgreSQL)
  - posts/: Posts and comments (PostgreSQL)
  - stats/: Analytics aggregation (ClickHouse, Kafka)
- tests/e2e/: End-to-end tests against a composed environment

Top-level Makefile provides common targets that wrap scripts/.

## Prerequisites

- Docker and Docker Compose
- make
- Go 1.22+ (only required to run unit tests locally without containers)

## Quick start

1) Build the Go builder image (used by service Dockerfiles):

```bash
make builder-image
```

2) Prepare environment variables. For a quick local run, copy test.env as a starting point:

```bash
cp test.env .env
```

Review and adjust .env paths and ports as needed. See the Environment variables section below.

3) Run the full stack:

```bash
make run
```

This will:
- Ensure data directories exist and have proper permissions (for Kafka)
- Build service images
- Start services as defined in deploy/docker-compose.yml

Once started, the Gateway listens on GATEWAY_SERVICE_PORT (default 8080) and exposes REST under /api/v1.

To stop the stack, press Ctrl+C in the compose terminal. If needed, run docker compose down with the same .env and compose file.

## Make targets

- make builder-image: Build the base Go builder image used by services
- make autogen: Generate protobuf stubs for all services
- make run: Start the full stack with docker compose using .env
- make unit-tests: Run unit tests for services
- make e2e-tests: Spin up stack using test.env, run e2e tests, then tear down
- make tests: Run both unit and e2e tests

## Environment variables

Values are loaded from .env by scripts/run-compose.sh and docker-compose.yml. The following variables are used:

- GATEWAY_SERVICE_PORT: Public HTTP port for the Gateway
- JWT_ED25519_PUBLIC_KEY: Hex-encoded public key used by Gateway and Posts to verify JWT
- JWT_ED25519_PRIVATE_KEY: Hex-encoded private key used by Accounts to issue JWT

- ACCOUNTS_SERVICE_PORT: gRPC port for Accounts service
- ACCOUNTS_POSTGRES_USER: PostgreSQL user for Accounts DB
- ACCOUNTS_POSTGRES_PASSWORD: PostgreSQL password for Accounts DB
- ACCOUNTS_POSTGRES_DATA: Host path for Accounts Postgres data volume

- POSTS_SERVICE_PORT: gRPC port for Posts service
- POSTS_POSTGRES_USER: PostgreSQL user for Posts DB
- POSTS_POSTGRES_PASSWORD: PostgreSQL password for Posts DB
- POSTS_POSTGRES_DATA: Host path for Posts Postgres data volume

- STATS_SERVICE_PORT: gRPC port for Stats service
- STATS_CLICKHOUSE_USER: ClickHouse user
- STATS_CLICKHOUSE_PASSWORD: ClickHouse password
- STATS_CLICKHOUSE_DATA: Host path for ClickHouse data volume

- STATS_KAFKA_DATA: Host path for Kafka data volume

Notes:
- The example keys in test.env are for local development only.
- scripts/run-compose.sh will chown STATS_KAFKA_DATA to uid/gid 1000 for Kafka container compatibility.

## Public API (Gateway)

Base URL: http://localhost:${GATEWAY_SERVICE_PORT}/api/v1

Selected endpoints (see services/gateway/internal/server/http_router.go for the full list):

- Auth and tokens:
  - POST /auth
  - POST /api_token

- Profiles:
  - POST /profile
  - GET /profile/:profile_id
  - PUT /profile/:profile_id
  - DELETE /profile/:profile_id

- Pages and posts:
  - GET /profile/:profile_id/page/settings
  - PUT /profile/:profile_id/page/settings
  - GET /profile/:profile_id/page/posts
  - POST /profile/:profile_id/page/posts

- Post operations:
  - GET /post/:post_id
  - PUT /post/:post_id
  - DELETE /post/:post_id
  - POST /post/:post_id/comments
  - GET /post/:post_id/comments
  - POST /post/:post_id/views
  - POST /post/:post_id/likes

- Metrics:
  - GET /post/:post_id/metric
  - GET /post/:post_id/metric_dynamics
  - GET /top10/posts
  - GET /top10/users

Request/response schemas are defined under services/gateway/api/ and swagger.yaml. The Gateway converts HTTP requests into internal gRPC calls.

## Development

- Code generation: make autogen to regenerate protobuf/gRPC stubs for all services.
- Unit tests: make unit-tests
- End-to-end tests: make e2e-tests (uses test.env and composes services in the background, then runs Go tests in tests/e2e).

## Troubleshooting

- Ports already in use: Adjust ports in .env.
- Permission issues on data directories: Ensure paths in .env exist and are writable. The scripts create and chown Kafka data directory automatically.
- e2e test failures: After make e2e-tests, check soa-e2e.log for composed service logs.

## Security

- JWT keys in test.env are only for local development. Provide your own secure keys in production.
- Do not commit real secrets. Prefer environment-specific .env files or secret stores.
