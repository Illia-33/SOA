all: generate build run

.PHONY:
build:
	go build -o gateway-service $(PWD)/cmd

.PHONY:
run: 
	$(PWD)/gateway-service
