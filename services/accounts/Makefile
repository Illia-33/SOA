all: proto build run

.PHONY: proto
proto:
	protoc $(PWD)/proto/*.proto \
	--proto_path=$(PWD)/proto \
	--go_out=$(PWD)/proto \
	--go-grpc_out=$(PWD)/proto

.PHONY: build
build: proto
	go build -o $(PWD)/accounts-service soa-socialnetwork/services/accounts/src

.PHONY: run
run: proto build
	$(PWD)/accounts-service

.PHONY: image
image:
	docker buildx build \
	--load \
	--file $(PWD)/Dockerfile \
	--tag soa-accounts-service \
	$(PWD)
