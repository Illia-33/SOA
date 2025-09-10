.PHONY: proto
proto:
	protoc $(PWD)/proto/*.proto \
	--proto_path=$(PWD)/proto \
	--go_out=$(PWD)/proto \
	--go-grpc_out=$(PWD)/proto

.PHONY: build
build:
	go build -o $(PWD)/stats-service $(PWD)/cmd

.PHONY: run
run: proto build
	$(PWD)/stats-service

.PHONY: image
image:
	docker buildx build \
	--load \
	--file $(PWD)/Dockerfile \
	--tag soa-stats-service \
	$(PWD)
