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

.PHONY: unit-tests
unit-tests:
	bash $(PWD)/scripts/run-unit-tests.sh

.PHONY: e2e-tests
e2e-tests: builder-image
	bash $(PWD)/scripts/run-e2e-tests.sh

.PHONY: tests
tests: unit-tests e2e-tests
