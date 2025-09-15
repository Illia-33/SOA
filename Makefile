all: builder-image

.PHONY: builder-image
builder-image:
	bash $(PWD)/scripts/build-builder-image.sh

.PHONY: autogen
autogen:
	bash $(PWD)/scripts/autogen.sh

.PHONY: run
run: builder-image
	bash $(PWD)/scripts/run-compose.sh

.PHONY: e2e-tests
e2e-tests: builder-image
	bash $(PWD)/scripts/run-e2e-tests.sh
