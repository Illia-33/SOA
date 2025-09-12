FROM soa-go-builder:latest AS builder

WORKDIR /source/services/stats

RUN make build

FROM ubuntu:24.04

WORKDIR /app

COPY --from=builder /source/services/stats/stats-service /app

CMD [ "/app/stats-service" ]
