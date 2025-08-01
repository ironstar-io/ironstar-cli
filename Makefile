VERSION_PATH ?= github.com/ironstar-io/ironstar-cli/internal/system/version
API_PATH     ?= github.com/ironstar-io/ironstar-cli/internal/api
BUILD_DATE   ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION      ?= $(shell git describe --tags)
GO_IMAGE     ?= golang:1.24.3

DOCKER_SCRIPT=docker run --rm \
		-v $(PWD)/.cache/go:/.cache \
		-v $(PWD):/src \
		-e GOCACHE=/.cache/go-build \
		-e GOMODCACHE=/.cache/go-mod \
		-e IRONSTAR_API_ADDRESS=https://nankai-dev:8443 \
		-e IRONSTAR_ARIMA_API_ADDRESS=http://arima:8000 \
		--network nankai_nankai \
		-w /src \
		-it \
		$(GO_IMAGE) \

GO_BUILD=go build \
	-trimpath \
	-ldflags "-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	-X $(API_PATH).version=$(VERSION)"

build:
	time \
	$(GO_BUILD) -o ./dist/iron

build-all: build-macos-amd64 build-macos-arm64 build-windows build-linux-amd64 build-linux-arm64

build-windows:
	env GOOS=windows GOARCH=amd64 \
	$(GO_BUILD) -o ./dist/iron-windows.exe

build-linux-amd64:
	env GOOS=linux GOARCH=amd64 \
	$(GO_BUILD) -o ./dist/iron-linux-amd64

build-linux-arm64:
	env GOOS=linux GOARCH=arm64 \
	$(GO_BUILD) -o ./dist/iron-linux-arm64

build-macos-amd64:
	env GOOS=darwin GOARCH=amd64 \
	$(GO_BUILD) -o ./dist/iron-macos

build-macos-arm64:
	env GOOS=darwin GOARCH=arm64 \
	$(GO_BUILD) -o ./dist/iron-macos-arm64

.PHONY: docker-run
docker-run: ## Run a CLI command in Docker, exiting immediately
docker-run:
	$(DOCKER_SCRIPT) /bin/bash -c "go run main.go $(CMD)"

.PHONY: docker-exec
docker-exec:
docker-exec:
	$(DOCKER_SCRIPT) /bin/bash

.PHONY: docker-test
docker-test:
docker-test:
	$(DOCKER_SCRIPT) go test ./...

clean:
	rm -rf ./dist/*

.PHONY: build build-windows build-linux-arm64 build-linux-amd-64 build-macos-amd64 build-macos-arm64 test clean
