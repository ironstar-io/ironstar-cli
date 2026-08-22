VERSION_PATH ?= github.com/ironstar-io/ironstar-cli/internal/system/version
API_PATH     ?= github.com/ironstar-io/ironstar-cli/internal/api
BUILD_DATE   ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION      ?= $(shell git describe --tags)
GO_IMAGE     ?= golang:1.27.0
GO_CACHE_DIR ?= $(PWD)/.cache/go
GOVULNCHECK_VERSION ?= v1.7.0

# Optional trust root for development networks that intercept HTTPS traffic:
#   make checks LOCAL_CA_CERT=/path/to/root-ca.crt
LOCAL_CA_CERT ?=
LOCAL_CA_MOUNT=$(if $(strip $(LOCAL_CA_CERT)),-v $(LOCAL_CA_CERT):/usr/local/share/ca-certificates/local-root.crt:ro)

DOCKER_BASE=docker run --rm \
		-v $(GO_CACHE_DIR):/.cache \
		-v $(PWD):/src \
		$(LOCAL_CA_MOUNT) \
		-e GOCACHE=/.cache/go-build \
		-e GOMODCACHE=/.cache/go-mod \
		-w /src

DOCKER_GO=$(DOCKER_BASE) $(GO_IMAGE)

DOCKER_NANKAI=$(DOCKER_BASE) \
		-e IRONSTAR_API_ADDRESS=https://nankai-dev:8443 \
		-e IRONSTAR_UPLOAD_DOMAIN=https://nankai-dev:8443 \
		--network nankai_nankai \
		-it \
		$(GO_IMAGE)

DOCKER_GO_RUN=$(DOCKER_GO) sh -c 'update-ca-certificates >/dev/null 2>&1 && exec "$$@"' sh

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

.PHONY: docker-run docker-exec docker-test fmt fmt/check vet test vuln build/check checks
docker-run: ## Run a CLI command in Docker, exiting immediately
	$(DOCKER_NANKAI) /bin/bash -c "update-ca-certificates >/dev/null 2>&1 && go run main.go $(CMD)"

docker-exec: ## Open an interactive Go shell on the Nankai Docker network
	$(DOCKER_NANKAI) /bin/bash -c "update-ca-certificates >/dev/null 2>&1 && exec /bin/bash"

docker-test: test ## Backwards-compatible alias for the Docker test target

fmt: ## Format Go source files in Docker
	$(DOCKER_GO) sh -c 'update-ca-certificates >/dev/null 2>&1 && gofmt -w main.go $$(find cmd internal -name "*.go" -type f)'

fmt/check: ## Fail when Go source files need formatting
	$(DOCKER_GO) sh -ec 'update-ca-certificates >/dev/null 2>&1; files="$$(gofmt -l main.go $$(find cmd internal -name "*.go" -type f))"; test -z "$$files" || { printf "Go files need formatting:\n%s\n" "$$files"; exit 1; }'

vet: ## Run go vet in Docker
	$(DOCKER_GO_RUN) go vet ./...

test: ## Run unit tests in Docker without requiring the Nankai network
	$(DOCKER_GO_RUN) go test ./...

vuln: ## Scan reachable Go code for known vulnerabilities
	$(DOCKER_GO_RUN) go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

build/check: ## Compile all packages in Docker without writing release binaries
	$(DOCKER_GO_RUN) go build ./...

checks: fmt/check vet test vuln build/check ## Run all local verification gates

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

.PHONY: build build-windows build-linux-arm64 build-linux-amd64 build-macos-amd64 build-macos-arm64 clean
