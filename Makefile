.DEFAULT_GOAL := init

init:
	@$(MAKE) buildclient

buildclient:
	@GOOS=js GOARCH=wasm go build -o ./client/wasm/main.wasm ./client/main.go && echo "WASM module built successfully"

example:
	@python3 -V > /dev/null 2>&1 || (echo "python3 is not installed. Please install it and try again." && exit 1) && \
		python3 -m http.server --directory ./client 9900

