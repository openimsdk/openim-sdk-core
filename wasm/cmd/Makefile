.PHONY: ios build install mac windows static/wasm_exec.js

WASM_BIN_NAME=openIM.wasm


# When adding a new target, remember to align the following command line with a tab at the beginning for proper indentation.
wasm:
	GOOS=js GOARCH=wasm go build -trimpath -ldflags "-s -w" -o ${WASM_BIN_NAME} main.go

build: static/main.wasm static/wasm_exec.js
	goexec 'http.ListenAndServe(`:9999`, http.FileServer(http.Dir(`.`)))'

static/wasm_exec.js:
	@mkdir -p static
	cp "$(shell go env GOROOT)/lib/wasm/wasm_exec.js" static

static/main.wasm : main.go
	GO111MODULE=auto GOOS=js GOARCH=wasm go build -o static/${WASM_BIN_NAME}.wasm main.go