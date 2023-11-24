package main

import (
	"context"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/msgtest/pressuser"
)

func main() {
	p := pressuser.NewPressureTester()
	f, r, err := p.SelectSample(10000, 0.01)
	if err != nil {
		log.ZError(context.Background(), "Sample UserID failed", err)
		return
	}
	if err := p.RegisterUsers(f, r, nil); err != nil {
		log.ZError(context.Background(), "Sample UserID failed", err)
		return
	}
	// init users
	p.InitUserConns(f, nil)
}
