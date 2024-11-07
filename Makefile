include .env

vendor:
	go mod tidy
	go mod vendor

build:
	go build -o bin/streamer main.go

run:
	@env $(cat $(PWD)/.env | xargs) bin/streamer

.PHONY: vendor build run
