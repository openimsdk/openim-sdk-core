.PHONY: ios build install mac windows

WASM_BIN_NAME=openIM.wasm


#新加target的时候，下面的命令行记得用tab顶格
wasm:
	GOOS=js GOARCH=wasm go build -trimpath -ldflags "-s -w" -o ${WASM_BIN_NAME} main.go

build: static/main.wasm static/wasm_exec.js
	goexec 'http.ListenAndServe(`:9999`, http.FileServer(http.Dir(`.`)))'

static/wasm_exec.js:
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" static

static/main.wasm : main.go
	GO111MODULE=auto GOOS=js GOARCH=wasm go build -o static/${WASM_BIN_NAME}.wasm main.go