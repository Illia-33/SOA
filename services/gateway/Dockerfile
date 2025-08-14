FROM soa-go-builder:latest AS builder

WORKDIR /source/services/accounts

RUN make proto

WORKDIR /source/services/gateway

RUN make build

FROM ubuntu:24.04

WORKDIR /app

COPY --from=builder /source/services/gateway/gateway-service /app

CMD [ "/app/gateway-service" ]
