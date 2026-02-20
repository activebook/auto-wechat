.PHONY: build run clean test fmt tidy release

BINARY_NAME := auto-wechat
MODULE     := $(shell go list -m)
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS    := -s -w \
              -X '$(MODULE)/internal.Version=$(VERSION)' \
              -X '$(MODULE)/internal.Commit=$(COMMIT)' \
              -X '$(MODULE)/internal.BuildDate=$(BUILD_DATE)'

OS   := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m | sed 's/x86_64/amd64/;s/arm64/arm64/')

DIST_DIR := dist

# ─── Development ─────────────────────────────────────────────────────────────

## build: compile the project
build:
	CGO_ENABLED=1 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) main.go

## run: build and execute the binary
run: build
	@if [ -f "./run.sh" ]; then ./run.sh $(ARGS); else ./$(BINARY_NAME) $(ARGS); fi

## test: run unit tests
test:
	go test ./...

## fmt: format the source code
fmt:
	go fmt ./...

## tag: create and push a git tag (usage: make tag v=v1.0.0)
tag:
	@if [ -z "$(v)" ]; then \
		LATEST=$$(git ls-remote --tags --refs origin | awk -F/ '{print $$NF}' | sort -V | tail -n 1); \
		if [ -z "$$LATEST" ]; then LATEST="v0.0.0"; fi; \
		echo "Latest remote tag: $$LATEST"; \
		NEXT=$$(echo $$LATEST | awk -F. '{print $$1"."$$2"."$$3+1}'); \
		echo "Suggested version: $$NEXT"; \
		echo "Usage: make tag v=$$NEXT"; \
		exit 0; \
	fi; \
	git tag $(v); \
	git push origin $(v)

# 
## upgrade: update dependencies to the latest versions
upgrade:
	go get -u ./...

## tidy: clean up and synchronize dependencies
tidy:
	go mod tidy

## clean: remove build artifacts and distributions
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -rf $(DIST_DIR)

# ─── Release ─────────────────────────────────────────────────────────────────

#RELEASE_NAME := $(BINARY_NAME)-$(OS)-$(ARCH)-$(VERSION)
RELEASE_NAME := $(BINARY_NAME)-$(OS)-$(ARCH)

## release: build a versioned binary for the current platform, bundle assets, and zip
release: clean
	@echo "→ Building $(BINARY_NAME) for $(OS)/$(ARCH)"
	@mkdir -p $(DIST_DIR)/$(RELEASE_NAME)/cfgs $(DIST_DIR)/$(RELEASE_NAME)/imgs
	CGO_ENABLED=1 go build \
		-trimpath \
		-ldflags "$(LDFLAGS)" \
		-o $(DIST_DIR)/$(RELEASE_NAME)/$(BINARY_NAME) \
		main.go
	cp cfgs/*.yml $(DIST_DIR)/$(RELEASE_NAME)/cfgs/
	cp imgs/*.png $(DIST_DIR)/$(RELEASE_NAME)/imgs/
	@if [ -f "run.sh" ]; then cp run.sh $(DIST_DIR)/$(RELEASE_NAME)/; fi
	cd $(DIST_DIR) && zip -r $(RELEASE_NAME).zip $(RELEASE_NAME)
	@echo "✓ Release archive: $(DIST_DIR)/$(RELEASE_NAME).zip"

## release-all: cross-compile for every supported platform (requires Docker / cross-toolchain)
release-all: clean
	@mkdir -p $(DIST_DIR)
	@echo "→ macOS arm64"
	@mkdir -p $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64-$(VERSION)/cfgs \
	           $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64-$(VERSION)/imgs
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build \
		-trimpath -ldflags "$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64-$(VERSION)/$(BINARY_NAME) main.go
	cp cfgs/*.yml $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64-$(VERSION)/cfgs/
	cp imgs/*.png $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64-$(VERSION)/imgs/
	cd $(DIST_DIR) && zip -r $(BINARY_NAME)-darwin-arm64-$(VERSION).zip $(BINARY_NAME)-darwin-arm64-$(VERSION)
	@echo "→ Windows amd64"
	@mkdir -p $(DIST_DIR)/$(BINARY_NAME)-windows-amd64-$(VERSION)/cfgs \
	           $(DIST_DIR)/$(BINARY_NAME)-windows-amd64-$(VERSION)/imgs
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build \
		-trimpath -ldflags "$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64-$(VERSION)/$(BINARY_NAME).exe main.go
	cp cfgs/*.yml $(DIST_DIR)/$(BINARY_NAME)-windows-amd64-$(VERSION)/cfgs/
	cp imgs/*.png $(DIST_DIR)/$(BINARY_NAME)-windows-amd64-$(VERSION)/imgs/
	cd $(DIST_DIR) && zip -r $(BINARY_NAME)-windows-amd64-$(VERSION).zip $(BINARY_NAME)-windows-amd64-$(VERSION)
	@echo "✓ Archives in $(DIST_DIR)/"
	@ls -lh $(DIST_DIR)/*.zip

## help: list available targets
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
