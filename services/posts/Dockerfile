FROM soa-go-builder:latest AS builder

WORKDIR /source/services/posts

RUN make build

FROM ubuntu:24.04

WORKDIR /app

COPY --from=builder /source/services/posts/posts-service /app

CMD [ "/app/posts-service" ]
