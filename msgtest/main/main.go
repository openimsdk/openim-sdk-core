package main

import (
	"context"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/msgtest/pressuser"
	"time"
)

func init() {

	if err := log.InitFromConfig("sdk.log", "sdk", 5,
		true, false, "./chat_log", 2, 24); err != nil {
		panic(err)
	}
}
func main() {
	ctx := context.Background()
	p := pressuser.NewPressureTester()
	f, r, err := p.SelectSample(10000, 0.01)
	if err != nil {
		log.ZError(ctx, "Sample UserID failed", err)
		return
	}
	log.ZDebug(ctx, "Sample UserID", "r", r)
	//if err := p.RegisterUsers(f, nil, nil); err != nil {
	//	log.ZError(ctx, "Sample UserID failed", err)
	//	return
	//}
	// init users
	p.InitUserConns(f, nil)
	time.Sleep(time.Hour * 60)

}
