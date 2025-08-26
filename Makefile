all: builder-image

.PHONY: builder-image
builder-image:
	docker buildx build \
	--load \
	--file $(PWD)/deploy/go-builder.Dockerfile \
	--build-context go-mods=$(PWD) \
	--tag soa-go-builder \
	$(PWD)

.PHONY: autogen
autogen:
	(cd $(PWD)/services/accounts; make proto)

.PHONY: run
run: builder-image
	docker compose \
	--env-file $(PWD)/.env \
	--file $(PWD)/deploy/docker-compose.yml \
	up \
	--build

.PHONY: e2e-tests
e2e-tests: builder-image
	docker compose --env-file $(PWD)/test.env --file $(PWD)/deploy/docker-compose.yml up --build --detach
	sleep 3s
	go test $(PWD)/tests/e2e || true
	docker compose --env-file $(PWD)/test.env --file $(PWD)/deploy/docker-compose.yml down
	sudo rm -r /temp/test-accounts-postgres
