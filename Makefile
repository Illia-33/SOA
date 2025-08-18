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