.PHONY: ios build install mac windows

WASM_BIN_NAME=openIM.wasm


#新加target的时候，下面的命令行记得用tab顶格
wasm:
	GOOS=js GOARCH=wasm go build -trimpath -ldflags "-s -w" -o ${WASM_BIN_NAME} main.go


