FROM soa-go-builder:latest AS builder

WORKDIR /source/services/accounts

COPY . .

RUN make proto
RUN make build

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /source/services/accounts/accounts-service ./

CMD [ "./accounts-service" ]
