package main

import (
	"context"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/msgtest/module"
	"time"
)

func init() {

	if err := log.InitFromConfig("sdk.log", "sdk", 5,
		true, false, "./", 2, 24); err != nil {
		panic(err)
	}
}
func main() {
	ctx := context.Background()
	p := module.NewPressureTester()
	f, r, err := p.SelectSample(1000, 0.01)
	if err != nil {
		log.ZError(ctx, "Sample UserID failed", err)
		return
	}
	log.ZDebug(ctx, "Sample UserID", "r", r, "length", len(f))
	time.Sleep(10 * time.Second)

	//if err := p.RegisterUsers(f, nil, nil); err != nil {
	//	log.ZError(ctx, "Sample UserID failed", err)
	//	return
	//}
	// init users
	p.InitUserConns(f)
	log.ZDebug(ctx, "all user init connect to server success,start send message")
	time.Sleep(10 * time.Second)
	p.SendSingleMessages(f, 10, time.Second)
	log.ZDebug(ctx, "message send finished start to check message")
	time.Sleep(20 * time.Second)
	p.CheckMsg()

	log.ZDebug(ctx, "message send finished start to check message")
	time.Sleep(time.Hour * 60)

}
