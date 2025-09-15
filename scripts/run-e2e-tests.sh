#!/bin/bash
set -e
source "$PWD/test.env"

sudo mkdir -p $STATS_KAFKA_DATA
sudo chown -R 1000:1000 $STATS_KAFKA_DATA

docker compose --env-file $PWD/test.env --file $PWD/deploy/docker-compose.yml up --build --detach

sleep 3s

go test -count=1 $PWD/tests/e2e || \
(echo "test failed, check soa-e2e.log" && docker compose --env-file $PWD/test.env --file $PWD/deploy/docker-compose.yml logs > $PWD/soa-e2e.log)

docker compose --env-file $PWD/test.env --file $PWD/deploy/docker-compose.yml down
sudo rm -r /temp/soa-e2e-test
