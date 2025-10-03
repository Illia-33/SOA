# Statistics Service

Aggregates and serves analytics for the SOA Social Network. Consumes domain events from Kafka and persists aggregated metrics to ClickHouse. Exposes a gRPC API for querying metrics and top lists.

- Language: Go
- Storage: ClickHouse
- Messaging: Kafka (topics: view, like, comment, registration, post)
- RPC: gRPC

## Responsibilities

- Consume events from Kafka (views, likes, comments, registrations, posts)
- Store raw and/or aggregated data in ClickHouse
- Provide metrics and dynamics for posts
- Provide top-10 posts and users endpoints

## gRPC API

Protocol buffers defined under services/stats/proto/stats_service.proto.
Generated code lives in the same directory.

## Environment variables

- STATS_SERVICE_PORT: gRPC port to listen on
- DB_HOST: ClickHouse host (e.g., stats-clickhouse)
- DB_PORT: ClickHouse port (e.g., 9000)
- DB_USER: ClickHouse username
- DB_PASSWORD: ClickHouse password
- KAFKA_HOST: Kafka broker host (e.g., stats-kafka)
- KAFKA_PORT: Kafka broker port (e.g., 9092)

## Run

Run with the full stack (recommended):

```bash
cp ../../test.env ../../.env   # from repo root
make -C ../../ run
```

Local run (requires ClickHouse and Kafka):

```bash
export STATS_SERVICE_PORT=50053
export DB_HOST=localhost; export DB_PORT=9000
export DB_USER=...; export DB_PASSWORD=...
export KAFKA_HOST=localhost; export KAFKA_PORT=9092

go run ./services/stats/cmd
```

## Development

- Regenerate protobuf stubs: `make proto`
- Unit tests for Stats: `make test`

## Directory structure

- cmd/: service entrypoint
- internal/
  - kafka/: consumer implementation
  - repo/: repositories for ClickHouse
  - server/: gRPC server setup
  - service/: business logic and converters
  - storage/clickhouse/: DB interactions
- db/migrations/: ClickHouse initialization
- proto/: protobuf definitions and generated code

## Notes

- Kafka topics are created by the init-stats-kafka container in docker-compose.
- Kafka UI is available at http://localhost:8088 by default when running via compose.
