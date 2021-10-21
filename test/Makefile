.PHONY: all build run gotool install clean help

BINARY_NAME=main
BIN_DIR=./
LAN_FILE=.go
GO_FILE:=${BINARY_NAME}${LAN_FILE}

all: gotool build

build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o  ${BINARY_NAME}  ${GO_FILE}

run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

install:
	make build
	mv main open_im_test_client

clean:
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi


