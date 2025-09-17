#!/bin/bash
set -e

source "$PWD/.env"

sudo mkdir -p $STATS_KAFKA_DATA
sudo chown -R 1000:1000 $STATS_KAFKA_DATA

docker compose \
    --env-file $PWD/.env \
    --file $PWD/deploy/docker-compose.yml \
    up \
    --build
