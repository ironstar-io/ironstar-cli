VERSION_PATH ?= gitlab.com/ironstar-io/ironstar-cli/internal/system/version
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= $(shell git describe --tags)

build:
	time \
	go build \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/ironstar

build-all: build-macos build-windows build-linux

build-windows:
	env GOOS=windows GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/ironstar-windows-amd64.exe

build-linux:
	env GOOS=linux GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/ironstar-linux-amd64

build-macos:
	env GOOS=darwin GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/ironstar-macos

usb-installer:
	cd installer && make build-macos
	cd installer && make build-windows
	cd installer && make build-linux
	cd installer && make build-docker-images
	make build-macos && cp -R ./dist/ironstar-macos ./installer/dist/ironstaraido/ironstar-macos
	make build-linux && cp -R ./dist/ironstar-linux-amd64 ./installer/dist/ironstaraido/ironstar-linux-amd64
	make build-windows && cp -R ./dist/ironstar-windows-amd64.exe ./installer/dist/ironstaraido/ironstar-windows-amd64.exe
	cp -R ./installer/README.md ./installer/dist/README.md

test:
	ginkgo test ./...

clean:
	rm -rf ./dist/*

.PHONY: build build-windows build-linux build-macos test clean