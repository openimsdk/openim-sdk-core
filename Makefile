# ==============================================================================
# define the default goal
#

SHELL := /bin/bash
DIRS=$(shell ls)
GO=go

# The NAME of the binary to build
NAME=ws_wrapper/cmd/open_im_sdk_server

# The directory to store the compiled binary
BIN_DIR=../../bin/

# The LAN_FILE extension for the source file
LAN_FILE=.go

# The full path to the source file
GO_FILE:=${NAME}${LAN_FILE}

# The default OS is Linux
OS:= $(or $(os),linux)

# The default architecture is amd64
ARCH:=$(or $(arch),amd64)

# Set the BINARY_NAME based on the OS
ifeq ($(OS),windows)
	BINARY_NAME=${NAME}.exe
else
	BINARY_NAME=${NAME}
endif

## build: Compile the binary
.PHONY: build
build:
	CGO_ENABLED=1 GOOS=${OS} GOARCH=${ARCH} go build -o ${BINARY_NAME}  ${GO_FILE}

## install: Install the binary to the BIN_DIR
.PHONY: install
install:build
	mv ${BINARY_NAME} ${BIN_DIR}

## clean: Clean the build artifacts
.PHONY: clean
clean:
	env GO111MODULE=on go clean -cache
	gomobile clean
	rm -fr build

## reset_remote_branch: Reset the remote branch
.PHONY: reset_remote_branch
reset_remote_branch:
	remote_branch=$(shell git rev-parse --abbrev-ref --symbolic-full-name @{u})
	git reset --hard $(remote_branch)
	git pull $(remote_branch)

## ios: Build the iOS framework
.PHONY: ios
ios:
	go get golang.org/x/mobile
	rm -rf build/ open_im_sdk/t_friend_sdk.go open_im_sdk/t_group_sdk.go  open_im_sdk/ws_wrapper/
	go mod download golang.org/x/exp
	GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" -o build/OpenIMCore.xcframework -target=ios ./open_im_sdk/ ./open_im_sdk_callback/

## android: Build the Android library
# Note: to build an AAR on Windows, gomobile, Android Studio, and the NDK must be installed.
# The NDK version tested by the OpenIM team was r20b.
# To build an AAR on Mac, gomobile, Android Studio, and the NDK version 20.0.5594570 must be installed.
.PHONY: android
android:
	go get golang.org/x/mobile/bind
	GOARCH=amd64 gomobile bind -v -trimpath -ldflags="-s -w" -o ./open_im_sdk.aar -target=android ./open_im_sdk/ ./open_im_sdk_callback/

## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\n\033[1mUsage: make <TARGETS> ...\033[0m\n\n\\033[1mTargets:\\033[0m\n\n"
	@sed -n 's/^##//p' $< | awk -F':' '{printf "\033[36m%-28s\033[0m %s\n", $$1, $$2}' | sed -e 's/^/ /'