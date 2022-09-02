VERSION_PATH ?= gitlab.com/ironstar-io/ironstar-cli/internal/system/version
BUILD_DATE   ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION      ?= $(shell git describe --tags)
GO_IMAGE     ?= golang:1.19

build:
	time \
	go build \
	-trimpath \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/iron

build-all: build-macos build-arm build-windows build-linux

build-windows:
	env GOOS=windows GOARCH=amd64 \
	go build \
	-trimpath \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/iron-windows.exe

build-linux:
	env GOOS=linux GOARCH=amd64 \
	go build \
	-trimpath \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/iron-linux-amd64

build-macos:
	env GOOS=darwin GOARCH=amd64 \
	go build \
	-trimpath \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/iron-macos

build-arm:
	env GOOS=darwin GOARCH=arm64 \
	go build \
	-trimpath \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/iron-macos-arm64


.PHONY: exec-docker
exec-docker: ## Run a CLI command in Docker
exec-docker:
	@docker run --rm \
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
		/bin/bash -c "go run main.go $(CMD)"

test:
	ginkgo test ./...

clean:
	rm -rf ./dist/*

.PHONY: build build-windows build-linux build-macos build-arm test clean