# layer8-genesis-repo
This repo contains the elementary component implementations to make the Layer8 system work: symmetric encryption, asymmetric key exchanges, OAuth flow, etc.

It was intended for the data transfer between clients and the Layer8 server to be over HTTP/2 using gRPC, but due to browser limitations, the data transfer will be over HTTP/1.1 using REST (at least for now until a way around [gRPC-Web](https://github.com/grpc/grpc-web) is found).

## Dependencies
- [Go](https://golang.org/doc/install)
- [Python 3](https://www.python.org/downloads/) (for the example client)
- [Make](https://www.gnu.org/software/make/)
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

## Setup
- Create a `.env` file in root directory and copy the contents of `.env.example` into it (or simply execute `cp .env.example .env`). Edit the values as needed.
- Execute `make setup` to install dependencies

## Build
- Execute `make build` to build the binaries (both WASM and the server binary)

## Run
- Execute `make run` to run the server (this runs both the REST and gRPC servers)
- To run only the REST server, execute `make run-rest`
- And to run only the gRPC server, execute `make run-grpc`

## Testing
- Execute `make test` to run unit tests
- An example client is provided in the `web` directory. To run it, execute `make example` and open the exposed URL in a browser

## Linting
- To run the linter, execute `make lint`

## Formatting
- Execute `make fmt` to run the formatter

## Flow Diagram
![Flow Diagram](./assets/flow_diagram.jpg)
