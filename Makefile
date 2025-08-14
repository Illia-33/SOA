all: builder-image

.PHONY: builder-image
builder-image:
	docker buildx build \
	--load \
	--file $(PWD)/docker/go-builder.Dockerfile \
	--build-context go-mods=$(PWD) \
	--tag soa-go-builder \
	$(PWD)