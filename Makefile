include .env

vendor:
	go mod tidy
	go mod vendor

build:
	go build -o bin/streamer main.go

run:
	bin/streamer run --config $(PWD)/resources/config.yaml

build.podman:
	@podman build . -t quay.io/ctupangiu/migration-event-streamer:latest
	@podman push quay.io/ctupangiu/migration-event-streamer:latest

infra.up:
	@podman play kube resources/dev.yml

infra.down:
	@podman kube down resources/dev.yml

.PHONY: vendor build run
