include .env
export

.DEFAULT_GOAL := setup
.PHONY: help

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "\033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

setup: ## Setup the project by installing dependencies and building the WASM module
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 && \
		export PATH=$$PATH:$(shell go env GOPATH)/bin && \
		$(MAKE) build

build: ## Build Go binary and WASM module
	@$(MAKE) build-go && $(MAKE) build-wasm

build-wasm: ## Build the WASM module
	@GOOS=js GOARCH=wasm go build -ldflags="-X main.ServerHost=$(SERVER_HOSTNAME) -X main.RESTPort=$(SERVER_REST_PORT) -X main.GRPCPort=$(SERVER_GRPC_PORT)" -o ./web/wasm/main.wasm ./web/main.go && echo "Go WASM built successfully"

build-go: ## Build the Go binary
	@go build -o ./bin/main ./cmd/main.go && echo "Go server built successfully"

example: ## Run the example client to test the WASM module
	@python3 -V > /dev/null 2>&1 || (echo "python3 is not installed. Please install it and try again." && exit 1) && \
		python3 -m http.server --directory ./web 9900

run: ## Run the server (REST API + gRPC)
	@test -f ./bin/main || (echo "Server binary not found. Please run 'make build' to build the server binary." && exit 1) && \
		./bin/main --rest-port=$(SERVER_REST_PORT) --grpc-port=$(SERVER_GRPC_PORT)

run-rest: ## Run the server (REST API)
	@test -f ./bin/main || (echo "Server binary not found. Please run 'make build' to build the server binary." && exit 1) && \
		./bin/main --rest-port=$(SERVER_REST_PORT)

run-grpc: ## Run the server (gRPC)
	@test -f ./bin/main || (echo "Server binary not found. Please run 'make build' to build the server binary." && exit 1) && \
		./bin/main --grpc-port=$(SERVER_GRPC_PORT)

test: ## Run tests
	@go test $(shell go list ./... | grep -v /web$)

lint: ## Run linter
	@golangci-lint run

fmt: ## Run gofmt
	@gofmt -w -s .

file = $(word 2, $(MAKECMDGOALS))
protoc: ## Compile a proto file to generate the gRPC code (e.g. make protoc file=grpc/proto/example.proto)
	@protoc --version > /dev/null 2>&1 || (echo "protoc is not installed. Please run 'make init' to install protoc." && exit 1) && \
		protoc -I ./api/grpc/proto --go_out=./api/grpc --go_opt=paths=source_relative --go-grpc_out=./api/grpc --go-grpc_opt=paths=source_relative $(file)
