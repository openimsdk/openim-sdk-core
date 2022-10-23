.PHONY: ios build install mac windows

WIN_BIN_NAME=localWs_win.exe
MAC_BIN_NAME=localWs_mac

#新加target的时候，下面的命令行记得用tab顶格
mac:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ${MAC_BIN_NAME} open_im_sdk_electron.go
windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ${WIN_BIN_NAME}  open_im_sdk_electron.go

