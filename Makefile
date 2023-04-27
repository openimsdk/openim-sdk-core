# ==============================================================================
# define the default goal
#

ROOT_PACKAGE=github.com/OpenIMSDK/Open-IM-SDK-Core

SHELL := /bin/bash
DIRS=$(shell ls)
GO=go

# include the common makefile
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# ROOT_DIR: root directory of the code base
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/. && pwd -P))
endif
# OUTPUT_DIR: The directory where the build output is stored.
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/bin
$(shell mkdir -p $(OUTPUT_DIR))
endif

ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --abbrev=0 --dirty --always --tags | sed 's/-/./g')
endif

# Check if the tree is dirty. default to dirty(maybe u should commit?)
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

# Define the directory you want to copyright
CODE_DIRS := $(ROOT_DIR) #$(ROOT_DIR)/pkg $(ROOT_DIR)/core $(ROOT_DIR)/integrationtest $(ROOT_DIR)/lib $(ROOT_DIR)/mock $(ROOT_DIR)/db $(ROOT_DIR)/openapi
FINDS := find $(CODE_DIRS)

ifndef V
MAKEFLAGS += --no-print-directory
endif

# Linux command settings
FIND := find . ! -path './image/*' ! -path './vendor/*' ! -path './bin/*'
XARGS := xargs -r
LICENSE_TEMPLATE ?= $(ROOT_DIR)/boilerplate.txt

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

BUILDFILE = "./main.go"
BUILDAPP = "$(OUTPUT_DIR)/"

.PHONY: ios build install android
.DEFAULT_GOAL := help

# ==============================================================================
# Targets

## all: Build all the necessary targets.
.PHONY: all
all: build

## build: Compile the binary
.PHONY: build
build:
	CGO_ENABLED=1 GOOS=${OS} GOARCH=${ARCH} go build -o ${BINARY_NAME}  ${GO_FILE}

## install: Install the binary to the BIN_DIR
.PHONY: install
install: build
	mv ${BINARY_NAME} ${BIN_DIR}

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

## tidy: tidy go.mod
.PHONY: tidy
tidy:
	@$(GO) mod tidy

## fmt: Run go fmt against code.
.PHONY: fmt
fmt:
	@$(GO) fmt ./...

## vet: Run go vet against code.
.PHONY: vet
vet:
	@$(GO) vet ./...

## generate: Run go generate against code.
.PHONY: generate
generate:
	@$(GO) generate ./...

## style: Code style -> fmt,vet,lint
.PHONY: style
style: fmt vet lint

## test: Run unit test
.PHONY: test
test: 
	@$(GO) test ./... 

## cover: Run unit test with coverage.
.PHONY: cover
cover: test
	@$(GO) test -cover

## lint: Run go lint against code.
.PHONY: lint
lint:
	@echo "===========> Run golangci to lint source codes"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

## copyright.verify: Validate boilerplate headers for assign files.
.PHONY: copyright-verify
copyright-verify: 
	@echo "===========> Validate boilerplate headers for assign files starting in the $(ROOT_DIR) directory"
	@addlicense -v -check -ignore **/test/** -f $(LICENSE_TEMPLATE) $(CODE_DIRS)
	@echo "===========> End of boilerplate headers check..."

## copyright-add: Add the boilerplate headers for all files.
.PHONY: copyright-add
copyright-add: 
	@echo "===========> Adding $(LICENSE_TEMPLATE) the boilerplate headers for all files"
	@addlicense -y $(shell date +"%Y") -v -c "OpenIM SDK." -f $(LICENSE_TEMPLATE) $(CODE_DIRS)
	@echo "===========> End the copyright is added..."

## clean: Clean the build artifacts
.PHONY: clean
clean:
	env GO111MODULE=on go clean -cache
	gomobile clean
	-rm -fvr build

## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\n\033[1mUsage: make <TARGETS> ...\033[0m\n\n\\033[1mTargets:\\033[0m\n\n"
	@sed -n 's/^##//p' $< | awk -F':' '{printf "\033[36m%-28s\033[0m %s\n", $$1, $$2}' | sed -e 's/^/ /'