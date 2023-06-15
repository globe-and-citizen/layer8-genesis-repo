# layer8-genesis-repo
This repo contains the elementary component implementations to make the Layer8 system work: symmetric encryption, asymmetric key exchanges, OAuth flow, etc.

It was intended for the data transfer between clients and the Layer8 server to be over HTTP/2 using gRPC, but due to browser limitations, the data transfer will be over HTTP/1.1 using REST (at least for now until a way around [gRPC-Web](https://github.com/grpc/grpc-web) is found).

## Dependencies
- [Go](https://golang.org/doc/install)
- [Python 3](https://www.python.org/downloads/) (for the example client)
- [Make](https://www.gnu.org/software/make/)
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

## Setup
- Execute `make setup` to install the dependencies

## Build
- Execute `make build` to build the binaries

## Run
- Execute `make run` to run the server

## Testing
- Execute `make test` to run the unit tests

## Linting
- Execute `make lint` to run the linter

## Formatting
- Execute `make fmt` to run the formatter

## Data Transfer Process
1. The client generates a public and private key, then shares the public key to the server for key exchange
2. The server generates a nonce and stores it in a database, then sends the nonce + its public key to the client
3. The client generates a shared secret using the its private key, server's public key and the nonce, then encrypts the data using the shared secret
4. The client sends the encrypted data to the server
5. The server looks up the nonce in the database, then generates a shared secret using its private key, the client's public key and the nonce
6. The server decrypts the data using the shared secret
7. The server sends the decrypted data to the internet
8. The server receives the response from the internet and encrypts it using the shared secret
9. The server sends the encrypted response to the client
10. The client decrypts the response using the shared secret
