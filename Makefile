.PHONY: ios build install

BINARY_NAME=ws_wrapper/cmd/open_im_sdk_server
BIN_DIR=../../bin/
LAN_FILE=.go
GO_FILE:=${BINARY_NAME}${LAN_FILE}

build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}  ${GO_FILE}
install:
	make build
	mv ${BINARY_NAME} ${BIN_DIR}
clean:
	env GO111MODULE=on go clean -cache
	gomobile clean
	rm -fr build

reset_remote_branch:
	remote_branch=$(shell git rev-parse --abbrev-ref --symbolic-full-name @{u})
	git reset --hard $(remote_branch)
	git pull $(remote_branch)

ios:
	go get golang.org/x/mobile
	rm -rf build/ open_im_sdk/t_friend_sdk.go open_im_sdk/t_group_sdk.go  open_im_sdk/ws_wrapper/
	go mod download golang.org/x/exp
	GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" -o build/OpenIMCore.xcframework -target=ios ./open_im_sdk/ ./open_im_sdk_callback/	
