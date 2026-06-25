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
		-e IRONSTAR_UPLOAD_DOMAIN=https://nankai-dev:8443 \
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

# --- Local signed macOS release -------------------------------------------
# Mirrors the CI release pipeline (.github/workflows/release.yml): build a
# universal (amd64+arm64) binary, code-sign it with the Developer ID under the
# hardened runtime, then notarise through App Store Connect. Only macOS is
# signed. CI imports the Developer ID from a secret into an ephemeral keychain;
# locally it is expected to already be in your keychain.
#
# Export the same env vars CI uses before running `make release-macos`:
#   APPLE_DEVELOPER_ID  Developer ID Application identity,
#                       e.g. "Developer ID Application: Ironstar ... (L7G23W3WF3)"
#   APPLE_TEAM_ID       Apple team id; the signature's TeamIdentifier is checked against it
#
# Notarisation runs only when all three of these are also set (otherwise it is
# skipped - a signed binary is enough to run locally from Terminal):
#   ASC_KEY_PATH        path to the App Store Connect API key (.p8)
#   ASC_KEY_ID          API key id
#   ASC_ISSUER_ID       API issuer id
#
# Override the embedded version with e.g. `make release-macos VERSION=v1.7.0-test`.
ENTITLEMENTS    ?= ./build/entitlements.plist
MACOS_UNIVERSAL ?= ./dist/iron-macos

.PHONY: build-macos-universal
build-macos-universal:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO_BUILD) -o ./dist/iron-macos-amd64
	env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO_BUILD) -o ./dist/iron-macos-arm64
	lipo -create -output $(MACOS_UNIVERSAL) ./dist/iron-macos-amd64 ./dist/iron-macos-arm64
	rm -f ./dist/iron-macos-amd64 ./dist/iron-macos-arm64
	lipo -info $(MACOS_UNIVERSAL)

.PHONY: sign-macos
sign-macos:
	@: $${APPLE_DEVELOPER_ID:?set APPLE_DEVELOPER_ID to your Developer ID Application identity}
	codesign --force --sign "$$APPLE_DEVELOPER_ID" --options runtime --timestamp --entitlements $(ENTITLEMENTS) $(MACOS_UNIVERSAL)
	codesign --verify --strict --verbose=2 $(MACOS_UNIVERSAL)
	@if [ -n "$$APPLE_TEAM_ID" ]; then \
		codesign -dv --verbose=4 $(MACOS_UNIVERSAL) 2>&1 | grep -q "TeamIdentifier=$$APPLE_TEAM_ID" \
			|| { echo "error: signed TeamIdentifier != APPLE_TEAM_ID ($$APPLE_TEAM_ID)"; exit 1; }; \
		echo "Verified TeamIdentifier=$$APPLE_TEAM_ID"; \
	fi

.PHONY: notarize-macos
notarize-macos:
	@if [ -z "$$ASC_KEY_PATH" ] || [ -z "$$ASC_KEY_ID" ] || [ -z "$$ASC_ISSUER_ID" ]; then \
		echo "ASC_KEY_PATH/ASC_KEY_ID/ASC_ISSUER_ID not all set - skipping notarisation."; \
		echo "(The binary is signed; that is enough to run it locally from Terminal.)"; \
	else \
		if [ ! -f "$$ASC_KEY_PATH" ]; then \
			echo "error: ASC_KEY_PATH ($$ASC_KEY_PATH) does not exist."; exit 1; \
		fi; \
		grep -q "BEGIN PRIVATE KEY" "$$ASC_KEY_PATH" 2>/dev/null || { \
			echo "error: ASC_KEY_PATH ($$ASC_KEY_PATH) is not an App Store Connect API key (.p8)."; \
			echo "       notarytool needs the .p8 PEM key from App Store Connect > Users and Access > Integrations > Keys."; \
			echo "       That is NOT the Developer ID .p12 used for code-signing (the .p12 backs APPLE_DEVELOPER_ID via your keychain)."; \
			exit 1; \
		}; \
		ZIP="$$(mktemp -d)/iron-macos.zip"; \
		/usr/bin/zip -j "$$ZIP" $(MACOS_UNIVERSAL) >/dev/null; \
		echo "Submitting $(MACOS_UNIVERSAL) to notarytool (waits for the result)..."; \
		xcrun notarytool submit "$$ZIP" --key "$$ASC_KEY_PATH" --key-id "$$ASC_KEY_ID" --issuer "$$ASC_ISSUER_ID" --wait; \
		st=$$?; rm -f "$$ZIP"; \
		[ $$st -eq 0 ] || { echo "notarisation failed; inspect with: xcrun notarytool log <id> --key \"$$ASC_KEY_PATH\" --key-id \"$$ASC_KEY_ID\" --issuer \"$$ASC_ISSUER_ID\""; exit $$st; }; \
	fi

# build -> sign -> (optional) notarise, in order even under `make -j`.
.PHONY: release-macos
release-macos:
	$(MAKE) build-macos-universal
	$(MAKE) sign-macos
	$(MAKE) notarize-macos
	@echo
	@echo "Signed macOS universal binary: $(MACOS_UNIVERSAL)"
	@$(MACOS_UNIVERSAL) version || true

clean:
	rm -rf ./dist/*

.PHONY: build build-windows build-linux-arm64 build-linux-amd-64 build-macos-amd64 build-macos-arm64 test clean
