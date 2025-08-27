FROM soa-go-builder:latest AS builder

WORKDIR /source/services/accounts

RUN make build

FROM ubuntu:24.04

WORKDIR /app

COPY --from=builder /source/services/accounts/accounts-service /app

CMD [ "/app/accounts-service" ]
