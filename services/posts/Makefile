.PHONY: proto
proto:
	protoc $(PWD)/proto/*.proto \
	--proto_path=$(PWD)/proto \
	--go_out=$(PWD)/proto \
	--go-grpc_out=$(PWD)/proto

.PHONY: build
build:
	go build -o $(PWD)/posts-service $(PWD)/cmd

.PHONY: run
run: proto build
	$(PWD)/posts-service

.PHONY: image
image:
	docker buildx build \
	--load \
	--file $(PWD)/Dockerfile \
	--tag soa-posts-service \
	$(PWD)
