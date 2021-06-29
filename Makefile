.PHONY: ios

clean:
	env GO111MODULE=on go clean -cache
	gomobile clean
	rm -fr build

reset_remote_branch:
	remote_branch=$(shell git rev-parse --abbrev-ref --symbolic-full-name @{u})
	git reset --hard $(remote_branch)
	git pull $(remote_branch)

ios: reset_remote_branch
	go get golang.org/x/mobile
	rm -rf build/ open_im_sdk/t_friend_sdk.go
	GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" -o build/OpenIMCore.framework -target=ios ./open_im_sdk/