package main

import (
	"open_im_sdk/pkg/utils"
	"open_im_sdk/testv3new"

	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
)

func main() {
	if err := log.InitFromConfig("sdk.log", "sdk", 6, true, false, "", 2, 24); err != nil {
		panic(err)
	}
	recvID := "9799811842"
	conversationNum := 10
	pressureTester := testv3new.NewPressureTester(testv3new.APIADDR, testv3new.WSADDR, testv3new.SECRET, testv3new.Admin)
	ctx := pressureTester.NewAdminCtx()
	ctx = mcontext.SetOperationID(ctx, utils.OperationIDGenerator())
	pressureTester.CreateConversations(ctx, conversationNum, recvID)
}
