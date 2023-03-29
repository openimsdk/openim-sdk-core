.PHONY: ios build install android

NAME=ws_wrapper/cmd/open_im_sdk_server
BIN_DIR=../../bin/
LAN_FILE=.go
GO_FILE:=${NAME}${LAN_FILE}
OS:= $(or $(os),linux)
ARCH:=$(or $(arch),amd64)

ifeq ($(OS),windows)

BINARY_NAME=${NAME}.exe

else

BINARY_NAME=${NAME}

endif


build:
	CGO_ENABLED=1 GOOS=${OS} GOARCH=${ARCH} go build -o ${BINARY_NAME}  ${GO_FILE}
install:build
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
#注：windows下打包成aar，保证gomobile,android studio以及NDK安装成功，NDK版本在window上官方测试为r20b,然后可以使用类似下面的命令生成aar
#   mac下打包成aar,保证gomobile,android studio以及NDK安装成功,NDK版本官方测试为20.0.5594570，使用如下命令生成aar
android:
	go get golang.org/x/mobile/bind
	GOARCH=amd64 gomobile bind -v -trimpath -ldflags="-s -w" -o ./open_im_sdk.aar -target=android ./open_im_sdk/ ./open_im_sdk_callback/
