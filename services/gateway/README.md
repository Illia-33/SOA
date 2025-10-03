# Gateway Service

Public HTTP entrypoint for the SOA Social Network. The Gateway exposes a REST API and delegates requests to internal services (Accounts, Posts, Stats) via gRPC. It also validates JWTs for protected operations.

- Language: Go
- Framework: Gin (HTTP)
- RPC: gRPC to internal services
- Auth: JWT (Ed25519)

## Responsibilities

- Expose REST API to clients under /api/v1
- Authenticate and authorize requests using JWT
- Translate HTTP requests to gRPC calls to internal services
- Compose responses and error handling for the public API

## API

OpenAPI schema: services/gateway/api/swagger.yaml

Key endpoints (see internal/server/http_router.go for full list):

- Auth and tokens:
  - POST /api/v1/auth
  - POST /api/v1/api_token

- Profiles:
  - POST /api/v1/profile
  - GET /api/v1/profile/:profile_id
  - PUT /api/v1/profile/:profile_id
  - DELETE /api/v1/profile/:profile_id

- Pages and posts:
  - GET /api/v1/profile/:profile_id/page/settings
  - PUT /api/v1/profile/:profile_id/page/settings
  - GET /api/v1/profile/:profile_id/page/posts
  - POST /api/v1/profile/:profile_id/page/posts

- Post operations:
  - GET /api/v1/post/:post_id
  - PUT /api/v1/post/:post_id
  - DELETE /api/v1/post/:post_id
  - POST /api/v1/post/:post_id/comments
  - GET /api/v1/post/:post_id/comments
  - POST /api/v1/post/:post_id/views
  - POST /api/v1/post/:post_id/likes

- Metrics:
  - GET /api/v1/post/:post_id/metric
  - GET /api/v1/post/:post_id/metric_dynamics
  - GET /api/v1/top10/posts
  - GET /api/v1/top10/users

Request/response schemas live in services/gateway/api/.

## Environment variables

- GATEWAY_SERVICE_PORT: Public HTTP port (e.g., 8080)
- JWT_ED25519_PUBLIC_KEY: Hex-encoded Ed25519 public key used to verify JWT
- ACCOUNTS_SERVICE_HOST: Hostname of Accounts service (gRPC)
- ACCOUNTS_SERVICE_PORT: Port of Accounts service (gRPC)
- POSTS_SERVICE_HOST: Hostname of Posts service (gRPC)
- POSTS_SERVICE_PORT: Port of Posts service (gRPC)
- STATS_SERVICE_HOST: Hostname of Stats service (gRPC)
- STATS_SERVICE_PORT: Port of Stats service (gRPC)

In Docker Compose these values are provided via .env and passed to the container.

## Run

Recommended: run with the full stack via the root Makefile:

```bash
cp test.env .env   # adjust as needed
make run
```

This builds and starts all services. The Gateway will listen on http://localhost:${GATEWAY_SERVICE_PORT}.

Local run (requires dependent services running):

```bash
export GATEWAY_SERVICE_PORT=8080
export JWT_ED25519_PUBLIC_KEY=...
export ACCOUNTS_SERVICE_HOST=localhost
export ACCOUNTS_SERVICE_PORT=50051
export POSTS_SERVICE_HOST=localhost
export POSTS_SERVICE_PORT=50052
export STATS_SERVICE_HOST=localhost
export STATS_SERVICE_PORT=50053

go run ./services/gateway/cmd
```

## Development

- Protobuf/gRPC stubs: generated from services/{accounts,posts,stats}/proto and used by the Gateway client stubs.
- Regenerate all protos (root): `make autogen`
- Unit tests (root): `make unit-tests`

## Directory structure

- cmd/: service entrypoint
- api/: HTTP request/response schemas and swagger
- internal/
  - server/: HTTP router and server wiring
  - service/: business logic and gRPC client invocations
  - query/: request parameter extraction and middleware
  - grpcutils/, httperr/: helpers

## Security notes

- Use strong, non-test JWT keys in production.
- Only provide minimal privileges for service-to-service networking.
