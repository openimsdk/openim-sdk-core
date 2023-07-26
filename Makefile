# ==============================================================================
# define the default goal
#

ROOT_PACKAGE=github.com/OpenIMSDK/Open-IM-SDK-Core

# Copyright 2023 OpenIM. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

###################################=> common commands <=#############################################
# ========================== Capture Environment ===============================
# get the repo root and output path
ROOT_PACKAGE=github.com/OpenIM/chat
OUT_DIR=$(REPO_ROOT)/_output
# ==============================================================================

# define the default goal
#

SHELL := /bin/bash
DIRS=$(shell ls)
GO=go

.DEFAULT_GOAL := help

# include the common makefile
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# ROOT_DIR: root directory of the code base
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/. && pwd -P))
endif
# OUTPUT_DIR: The directory where the build output is stored.
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif

# BIN_DIR: The directory where the build output is stored.
ifeq ($(origin BIN_DIR),undefined)
BIN_DIR := $(OUTPUT_DIR)/bin
$(shell mkdir -p $(BIN_DIR))
endif

ifeq ($(origin TOOLS_DIR),undefined)
TOOLS_DIR := $(OUTPUT_DIR)/tools
$(shell mkdir -p $(TOOLS_DIR))
endif

ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
$(shell mkdir -p $(TMP_DIR))
endif

ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match="v*" --dirty | sed 's/-/./g')	#v2.3.3.631.g00abdc9b.dirty
endif

# Check if the tree is dirty. default to dirty(maybe u should commit?)
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

IMG ?= openim_chat:latest

BUILDFILE = "./main.go"
BUILDAPP = "$(OUTPUT_DIR)/"

# Define the directory you want to copyright
CODE_DIRS := $(ROOT_DIR)/ #$(ROOT_DIR)/pkg $(ROOT_DIR)/core $(ROOT_DIR)/integrationtest $(ROOT_DIR)/lib $(ROOT_DIR)/mock $(ROOT_DIR)/db $(ROOT_DIR)/openapi
FINDS := find $(CODE_DIRS)

ifndef V
MAKEFLAGS += --no-print-directory
endif

# The OS must be linux when building docker images
# !WARNING: linux_mips64 linux_mips64le
PLATFORMS ?= linux_s390x darwin_amd64 windows_amd64 linux_amd64 linux_arm64 linux_ppc64le
# The OS can be linux/windows/darwin when building binaries
# PLATFORMS ?= darwin_amd64 windows_amd64 linux_amd64 linux_arm64

# Set a specific PLATFORM
ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
	# Use linux as the default OS when building images
	IMAGE_PLAT := linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif

# Copy githook scripts when execute makefile
# TODO! GIT_FILE_SIZE_LIMIT=42000000 git commit -m "This commit is allowed file sizes up to 42MB"
COPY_GITHOOK:=$(shell cp -f scripts/githooks/* .git/hooks/; chmod +x .git/hooks/*)

# Linux command settings
FIND := find . ! -path './image/*' ! -path './vendor/*' ! -path './bin/*'
XARGS := xargs -r

# ==============================================================================
# TODO: License selection
# LICENSE_TEMPLATE ?= $(ROOT_DIR)/scripts/LICENSE/license_templates.txt	# MIT License
LICENSE_TEMPLATE ?= $(ROOT_DIR)/scripts/LICENSE/LICENSE_TEMPLATES  # Apache License

# COMMA: Concatenate multiple strings to form a list of strings
COMMA := ,
# SPACE: Used to separate strings
SPACE :=
# SPACE: Replace multiple consecutive Spaces with a single space
SPACE +=

# ==============================================================================
# Build definition

GO_SUPPORTED_VERSIONS ?= 1.18|1.19|1.20
GO_LDFLAGS += -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
ifneq ($(DLV),)
	GO_BUILD_FLAGS += -gcflags "all=-N -l"
	LDFLAGS = ""
endif
GO_BUILD_FLAGS += -ldflags "$(GO_LDFLAGS)"

ifeq ($(GOOS),windows)
	GO_OUT_EXT := .exe
endif

ifeq ($(ROOT_PACKAGE),)
	$(error the variable ROOT_PACKAGE must be set prior to including golang.mk)
endif

GOPATH := $(shell go env GOPATH)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif

COMMANDS ?= $(filter-out %.md, $(wildcard ${ROOT_DIR}/cmd/*))
BINS ?= $(foreach cmd,${COMMANDS},$(notdir ${cmd}))

ifeq (${COMMANDS},)
  $(error Could not determine COMMANDS, set ROOT_DIR or run in source dir)
endif
ifeq (${BINS},)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif

EXCLUDE_TESTS=github.com/OpenIMSDK/chat/test

# ==============================================================================
# Build

## all: Build all the necessary targets.
.PHONY: all
all: copyright-verify build # tidy lint cover

## build: Build binaries by default.
.PHONY: build
build:
	@CGO_ENABLED=1 GOOS=${OS} GOARCH=${ARCH} go build -o ${BIN_DIR} ./cmd/main.go
# build: go.build.verify $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS)))

.PHONY: build.%
build.%:
	@echo "$(shell go version)"
	@echo "===========> Building binary $(BUILDAPP) *[Git Info]: $(VERSION)-$(GIT_COMMIT)"
	@export CGO_ENABLED=0 && GOOS=linux go build -o $(BUILDAPP)/$*/ -ldflags '-s -w' $*/example/$(BUILDFILE)

.PHONY: go.build.verify
go.build.verify:
ifneq ($(shell $(GO) version | grep -q -E '\bgo($(GO_SUPPORTED_VERSIONS))\b' && echo 0 || echo 1), 0)
	$(error unsupported go version. Please make install one of the following supported version: '$(GO_SUPPORTED_VERSIONS)')
endif

## go.build: Build the binary file of the specified platform.
.PHONY: go.build.%
go.build.%:
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "=====> COMMAND=$(COMMAND)"
	@echo "=====> PLATFORM=$(PLATFORM)"
	@echo "=====> BIN_DIR=$(BIN_DIR)"
	@echo "===========> Building binary $(COMMAND) $(VERSION) for $(OS)_$(ARCH)"
	@mkdir -p $(BIN_DIR)/platforms/$(OS)/$(ARCH)
	@if [ "$(COMMAND)" = "rpc" ] || [ "$(COMMAND)" = "api" ]; then \
		for d in $(wildcard $(ROOT_DIR)/cmd/$(COMMAND)/*/); do \
			cd $$d && CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o \
			$(BIN_DIR)/platforms/$(OS)/$(ARCH)/$$(basename $$d)$(GO_OUT_EXT) .; \
		done; \
	else \
		@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) -o \
		$(BIN_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) $(ROOT_PACKAGE)/cmd/$(COMMAND)/main.go; \
	fi

## build-multiarch: Build binaries for multiple platforms.
.PHONY: build-multiarch
build-multiarch: go.build.verify $(foreach p,$(PLATFORMS),$(addprefix go.build., $(addprefix $(p)., $(BINS))))
# ==============================================================================
# Targets

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

# Targets
.PHONY: release
release: release.verify release.ensure-tag
	@scripts/release.sh

.PHONY: install.gsemver
release.verify: install.git-chglog install.github-release install.coscmd

.PHONY: release.tag
release.tag: install.gsemver release.ensure-tag
	@git push origin `git describe --tags --abbrev=0`

.PHONY: release.ensure-tag
release.ensure-tag: install.gsemver
	@scripts/ensure_tag.sh

## tidy: tidy go.mod
.PHONY: tidy
tidy:
	@$(GO) mod tidy

## style: Code style -> fmt,vet,lint
.PHONY: style
style: fmt vet lint

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

## lint: Run go lint against code.
.PHONY: lint
lint: tools.verify.golangci-lint
	@echo "===========> Run golangci to lint source codes"
	@$(TOOLS_DIR)/golangci-lint run -c $(ROOT_DIR)/.golangci.yml $(ROOT_DIR)/...

## test: Run unit test
.PHONY: test
test: 
	@$(GO) test ./... 

## cover: Run unit test with coverage.
.PHONY: cover
cover: test
	@$(GO) test -cover

## docker-build: Build docker image with the manager.
.PHONY: docker-build
docker-build:
	docker build -t ${IMG} .

## docker-push: Push docker image with the manager.
.PHONY: docker-push
docker-push:
	docker push ${IMG}

## docker-buildx-push: Push docker image with the manager using buildx.
.PHONY: docker-buildx-push
docker-buildx-push:
	docker buildx build --platform linux/arm64,linux/amd64 -t ${IMG} . --push

## copyright-verify: Validate boilerplate headers for assign files.
.PHONY: copyright-verify
copyright-verify: tools.verify.addlicense copyright-add
	@echo "===========> Validate boilerplate headers for assign files starting in the $(ROOT_DIR) directory"
	@$(TOOLS_DIR)/addlicense -v -check -ignore **/test/** -f $(LICENSE_TEMPLATE) $(CODE_DIRS)
	@echo "===========> End of boilerplate headers check..."

## copyright-add: Add the boilerplate headers for all files.
.PHONY: copyright-add
copyright-add: tools.verify.addlicense
	@echo "===========> Adding $(LICENSE_TEMPLATE) the boilerplate headers for all files"
	@$(TOOLS_DIR)/addlicense -y $(shell date +"%Y") -v -c "OpenIM open source community." -f $(LICENSE_TEMPLATE) $(CODE_DIRS)
	@echo "===========> End the copyright is added..."

## clean: Clean all builds.
.PHONY: clean
clean:
	@echo "===========> Cleaning all builds TMP_DIR($(TMP_DIR)) AND BIN_DIR($(BIN_DIR))"
	@-rm -vrf $(TMP_DIR) $(BIN_DIR)
	@echo "===========> End clean..."

## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\n\033[1mUsage: make <TARGETS> ...\033[0m\n\n\\033[1mTargets:\\033[0m\n\n"
	@sed -n 's/^##//p' $< | awk -F':' '{printf "\033[36m%-28s\033[0m %s\n", $$1, $$2}' | sed -e 's/^/ /'
	
######################################=> common tools<= ############################################
# tools

BUILD_TOOLS ?= go-gitlint golangci-lint goimports addlicense deepcopy-gen conversion-gen ginkgo go-junit-report 

## tools.verify.%: Check if a tool is installed and install it
.PHONY: tools.verify.%
tools.verify.%:
	@echo "===========> Verifying $* is installed"
	@if [ ! -f $(TOOLS_DIR)/$* ]; then GOBIN=$(TOOLS_DIR) $(MAKE) tools.install.$*; fi
	@echo "===========> $* is install in $(TOOLS_DIR)/$*"

# tools: Install a must tools
.PHONY: tools
tools: $(addprefix tools.verify., $(BUILD_TOOLS))

# tools.install.%: Install a single tool in $GOBIN/
.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $,The default installation path is $(GOBIN)/$*"
	@$(MAKE) install.$*

.PHONY: install.golangci-lint
install.golangci-lint:
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: install.goimports
install.goimports:
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

.PHONY: install.addlicense
install.addlicense:
	@$(GO) install github.com/google/addlicense@latest

.PHONY: install.deepcopy-gen
install.deepcopy-gen:
	@$(GO) install k8s.io/code-generator/cmd/deepcopy-gen@latest

.PHONY: install.conversion-gen
install.conversion-gen:
	@$(GO) install k8s.io/code-generator/cmd/conversion-gen@latest

.PHONY: install.ginkgo
install.ginkgo:
	@$(GO) install github.com/onsi/ginkgo/ginkgo@v1.16.2

.PHONY: install.go-gitlint
# wget -P _output/tools/ https://openim-1306374445.cos.ap-guangzhou.myqcloud.com/openim/tools/go-gitlint
# go install github.com/antham/go-gitlint/cmd/gitlint@latest
install.go-gitlint:
	@wget -q https://openim-1306374445.cos.ap-guangzhou.myqcloud.com/openim/tools/go-gitlint -O ${TOOLS_DIR}/go-gitlint
	@chmod +x ${TOOLS_DIR}/go-gitlint

.PHONY: install.go-junit-report
install.go-junit-report:
	@$(GO) install github.com/jstemmer/go-junit-report@latest

# ==============================================================================
# Tools that might be used include go gvm, cos
#

## install.kube-score: Install kube-score, used to check kubernetes yaml files
.PHONY: install.kube-score
install.kube-score:
	@$(GO) install github.com/zegl/kube-score/cmd/kube-score@latest

## install.kubeconform: Install kubeconform, used to check kubernetes yaml files
.PHONY: install.kubeconform
install.kubeconform:
	@$(GO) install github.com/yannh/kubeconform/cmd/kubeconform@latest

## install.gsemver: Install gsemver, used to generate semver
.PHONY: install.gsemver
install.gsemver:
	@$(GO) install github.com/arnaud-deprez/gsemver@latest

## install.git-chglog: Install git-chglog, used to generate changelog
.PHONY: install.git-chglog
install.git-chglog:
	@$(GO) install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

## install.github-release: Install github-release, used to create github release
.PHONY: install.github-release
install.github-release:
	@$(GO) install github.com/github-release/github-release@latest

## install.coscli: Install coscli, used to upload files to cos
# example: ./coscli  cp/sync -r /root/workspaces/kubecub/chat/ cos://kubecub-1306374445/code/ -e cos.ap-hongkong.myqcloud.com
# https://cloud.tencent.com/document/product/436/71763
# kubecub/*
# - code/
# - docs/
# - images/
# - scripts/
.PHONY: install.coscli
install.coscli:
	@wget -q https://github.com/tencentyun/coscli/releases/download/v0.13.0-beta/coscli-linux -O ${TOOLS_DIR}/coscli
	@chmod +x ${TOOLS_DIR}/coscli

## install.coscmd: Install coscmd, used to upload files to cos
.PHONY: install.coscmd
install.coscmd:
	@if which pip &>/dev/null; then pip install coscmd; else pip3 install coscmd; fi

## install.delve: Install delve, used to debug go program
.PHONY: install.delve
install.delve:
	@$(GO) install github.com/go-delve/delve/cmd/dlv@latest

## install.air: Install air, used to hot reload go program
.PHONY: install.air
install.air:
	@$(GO) install github.com/cosmtrek/air@latest

## install.gvm: Install gvm, gvm is a Go version manager, built on top of the official go tool.
.PHONY: install.gvm
install.gvm:
	@echo "===========> Installing gvm,The default installation path is ~/.gvm/scripts/gvm"
	@bash < <(curl -s -S -L https://raw.gitee.com/moovweb/gvm/master/binscripts/gvm-installer)
	@$(shell source /root/.gvm/scripts/gvm)

## install.golines: Install golines, used to format long lines
.PHONY: install.golines
install.golines:
	@$(GO) install github.com/segmentio/golines@latest

## install.go-mod-outdated: Install go-mod-outdated, used to check outdated dependencies
.PHONY: install.go-mod-outdated
install.go-mod-outdated:
	@$(GO) install github.com/psampaz/go-mod-outdated@latest

## install.mockgen: Install mockgen, used to generate mock functions
.PHONY: install.mockgen
install.mockgen:
	@$(GO) install github.com/golang/mock/mockgen@latest

## install.gotests: Install gotests, used to generate test functions
.PHONY: install.gotests
install.gotests:
	@$(GO) install github.com/cweill/gotests/gotests@latest

## install.protoc-gen-go: Install protoc-gen-go, used to generate go source files from protobuf files
.PHONY: install.protoc-gen-go
install.protoc-gen-go:
	@$(GO) install github.com/golang/protobuf/protoc-gen-go@latest

## install.cfssl: Install cfssl, used to generate certificates
.PHONY: install.cfssl
install.cfssl:
	@$(ROOT_DIR)/scripts/install/install.sh OpenIM::install::install_cfssl

## install.depth: Install depth, used to check dependency tree
.PHONY: install.depth
install.depth:
	@$(GO) install github.com/KyleBanks/depth/cmd/depth@latest

## install.go-callvis: Install go-callvis, used to visualize call graph
.PHONY: install.go-callvis
install.go-callvis:
	@$(GO) install github.com/ofabry/go-callvis@latest

## install.gothanks: Install gothanks, used to thank go dependencies
.PHONY: install.gothanks
install.gothanks:
	@$(GO) install github.com/psampaz/gothanks@latest

## install.richgo: Install richgo
.PHONY: install.richgo
install.richgo:
	@$(GO) install github.com/kyoh86/richgo@latest

## install.rts: Install rts
.PHONY: install.rts
install.rts:
	@$(GO) install github.com/galeone/rts/cmd/rts@latest