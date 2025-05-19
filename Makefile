# Makefile for groq-chat CLI
.PHONY: all build clean test release

# Binary name as a variable
APP_NAME = groq-chat

RELEASE_DIR := bin/release
FLAT_DIR := $(RELEASE_DIR)/flat

# Version number
VERSION = v1.0.0
#VERSION ?= $(shell git describe --tags --always --dirty)

# Command prefix for silent/verbose builds (use VERBOSE=1 for verbose output)
V = @
ifeq ($(VERBOSE),1)
    V =
endif

# Detect Windows for .exe suffix in build
ifeq ($(OS),Windows_NT)
    APP_SUFFIX = .exe
else
    APP_SUFFIX =
endif

all: build

build: deps
	$(V)mkdir -p bin
	$(V)go build -ldflags "-X main.version=$(VERSION)" -o bin/$(APP_NAME)$(APP_SUFFIX) ./cmd/groq-cli-chat
# UPX disabled due to lack of official v5 binary for macOS; uncomment to enable
#	$(V)command -v upx >/dev/null && upx --best bin/$(APP_NAME)$(APP_SUFFIX) || echo "UPX not found, skipping compression for build"

clean:
	$(V)rm -rf bin

test:
	$(V)go test ./...

deps:
	$(V)go mod tidy
	$(V)go mod download

release: clean release-artifacts release-flatten

release-artifacts: deps
	$(V)mkdir -p bin/release/linux-amd64 bin/release/linux-arm64 bin/release/macos-amd64 bin/release/macos-arm64 bin/release/windows-amd64 bin/release/windows-arm64
	$(V)CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/release/linux-amd64/$(APP_NAME) ./cmd/groq-cli-chat
	$(V)CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/release/linux-arm64/$(APP_NAME) ./cmd/groq-cli-chat
	$(V)CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/release/macos-amd64/$(APP_NAME) ./cmd/groq-cli-chat
	$(V)CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/release/macos-arm64/$(APP_NAME) ./cmd/groq-cli-chat
	$(V)CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/release/windows-amd64/$(APP_NAME).exe ./cmd/groq-cli-chat
	$(V)CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/release/windows-arm64/$(APP_NAME).exe ./cmd/groq-cli-chat

	$(V)tar -czf bin/release/$(APP_NAME)-linux-amd64.tar.gz -C bin/release/linux-amd64 $(APP_NAME)
	$(V)tar -czf bin/release/$(APP_NAME)-linux-arm64.tar.gz -C bin/release/linux-arm64 $(APP_NAME)
	$(V)tar -czf bin/release/$(APP_NAME)-macos-amd64.tar.gz -C bin/release/macos-amd64 $(APP_NAME)
	$(V)tar -czf bin/release/$(APP_NAME)-macos-arm64.tar.gz -C bin/release/macos-arm64 $(APP_NAME)

	$(V)cd bin/release/windows-amd64 && zip -q ../$(APP_NAME)-windows-amd64.zip $(APP_NAME).exe || echo "zip command failed"
	$(V)cd bin/release/windows-arm64 && zip -q ../$(APP_NAME)-windows-arm64.zip $(APP_NAME).exe || echo "zip command failed"

release-flatten:
	$(V)mkdir -p $(FLAT_DIR)
	@echo "Flattening and renaming release artifacts..."
	@for file in $(RELEASE_DIR)/$(APP_NAME)-* ; do \
		[ -f "$$file" ] || continue ; \
		platform_arch_ext=$$(echo "$${file##*/}" | sed -E 's/^$(APP_NAME)-//') ; \
		cp "$$file" "$(FLAT_DIR)/$$platform_arch_ext" ; \
		echo "Copied: $$file -> $(FLAT_DIR)/$$platform_arch_ext" ; \
	done

# Docker image name and tag
DOCKER_REPO = oleksiyml/groq-chat
#TAG = v$(VERSION)
TAG = $(VERSION)

docker-multiarch:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(DOCKER_REPO):$(TAG) \
		-t $(DOCKER_REPO):latest \
		--file Dockerfile.multiarch \
		--push .