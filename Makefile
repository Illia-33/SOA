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
	(cd $(PWD)/services/posts; make proto)
	(cd $(PWD)/services/stats; make proto)

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
	go test -count=1 $(PWD)/tests/e2e || \
	(echo "test failed, check soa-e2e.log" && docker compose --env-file $(PWD)/test.env --file $(PWD)/deploy/docker-compose.yml logs > $(PWD)/soa-e2e.log)
	docker compose --env-file $(PWD)/test.env --file $(PWD)/deploy/docker-compose.yml down
	sudo rm -r /temp/test-accounts-postgres
	sudo rm -r /temp/test-posts-postgres
