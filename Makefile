GOLANGCI_LINT_CACHE?=/tmp/practicum-golangci-lint-cache

.PHONY: golangci-lint
.PHONY: golangci-lint-format
.PHONY: build-gophkeeper
.PHONY: build-gophkeeper-client
.PHONY: build-gophkeeper-proto
.PHONY: compile-all
.PHONY: run-server
.PHONY: run-client

# ====================================== VARIABLES =====================================

# build variables
ifdef GOPHKEEPER_BUILD_INFO_VERSION
GOPHKEEPER_BUILD_INFO_VERSION := $(GOPHKEEPER_BUILD_INFO_VERSION)
else
GOPHKEEPER_BUILD_INFO_VERSION := 0.0.1
endif

ifdef GOPHKEEPER_BUILD_INFO_DATE
GOPHKEEPER_BUILD_INFO_DATE := $(GOPHKEEPER_BUILD_INFO_DATE)
else
GOPHKEEPER_BUILD_INFO_DATE := $(shell date +'%Y-%m-%d_%H:%M:%S')
endif

ifdef GOPHKEEPER_BUILD_INFO_COMMIT_HASH
GOPHKEEPER_BUILD_INFO_COMMIT_HASH := $(GOPHKEEPER_BUILD_INFO_COMMIT_HASH)
else
GOPHKEEPER_BUILD_INFO_COMMIT_HASH := $(shell git rev-parse --short HEAD)
endif

LDFLAGS = -ldflags "-X main.buildInfoVersion=$(GOPHKEEPER_BUILD_INFO_VERSION) -X main.buildInfoDate=$(GOPHKEEPER_BUILD_INFO_DATE) -X main.buildInfoCommitHash=$(GOPHKEEPER_BUILD_INFO_COMMIT_HASH)"

# Directories
BIN_DIR = bin

# ====================================== COMMANDS =====================================

golangci-lint:
	mkdir -p ./golangci-lint
	golangci-lint run -c .golangci.yml

golangci-lint-format:
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json
	rm ./golangci-lint/report-unformatted.json

build-gophkeeper:
	go build -o ./cmd/gophkeeper/gophkeeper -ldflags $(LDFLAGS) ./cmd/gophkeeper/

build-gophkeeper-client:
	go build -o ./cmd/gophkeeper/gophkeeper-client -ldflags $(LDFLAGS) ./cmd/gophkeeper-client/

run-server:
	go run ./cmd/gophkeeper/

run-client:
	go run ./cmd/gophkeeper-client/main.go

build-gophkeeper-proto:
	protoc --go_out=.\
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
    --go_opt=default_api_level=API_OPAQUE \
    api/gophkeeper/secret.proto
	protoc --go_out=.\
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
    --go_opt=default_api_level=API_OPAQUE \
    api/gophkeeper/auth.proto


compile-all:
	@mkdir -p $(BIN_DIR)
	@mkdir -p $(BIN_DIR)/$(GOPHKEEPER_BUILD_INFO_VERSION)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(GOPHKEEPER_BUILD_INFO_VERSION)/gophkeeper-$(GOPHKEEPER_BUILD_INFO_VERSION)-linux-amd64 ./cmd/gophkeeper
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(GOPHKEEPER_BUILD_INFO_VERSION)/gophkeeper-$(GOPHKEEPER_BUILD_INFO_VERSION)-darwin-amd64 ./cmd/gophkeeper
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/$(GOPHKEEPER_BUILD_INFO_VERSION)/gophkeeper-$(GOPHKEEPER_BUILD_INFO_VERSION)-darwin-arm64 ./cmd/gophkeeper
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(GOPHKEEPER_BUILD_INFO_VERSION)/gophkeeper-$(GOPHKEEPER_BUILD_INFO_VERSION)-windows-amd64.exe ./cmd/gophkeeper
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(GOPHKEEPER_BUILD_INFO_VERSION)/gophkeeper-client-$(GOPHKEEPER_BUILD_INFO_VERSION)-linux-amd64 ./cmd/gophkeeper-client
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(GOPHKEEPER_BUILD_INFO_VERSION)/gophkeeper-client-$(GOPHKEEPER_BUILD_INFO_VERSION)-darwin-amd64 ./cmd/gophkeeper-client
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/$(GOPHKEEPER_BUILD_INFO_VERSION)/gophkeeper-client-$(GOPHKEEPER_BUILD_INFO_VERSION)-darwin-arm64 ./cmd/gophkeeper-client
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/$(GOPHKEEPER_BUILD_INFO_VERSION)/gophkeeper-client-$(GOPHKEEPER_BUILD_INFO_VERSION)-windows-amd64.exe ./cmd/gophkeeper-client
