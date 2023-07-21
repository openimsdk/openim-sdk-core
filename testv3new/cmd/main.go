package main

import (
	"open_im_sdk/pkg/utils"
	"open_im_sdk/testv3new"
	"open_im_sdk/testv3new/testcore"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
)

func main() {
	if err := log.InitFromConfig("sdk.log", "sdk", 3, true, false, "", 2); err != nil {
		panic(err)
	}
	userID := "4844258055"
	recvID := "4950983283"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI0ODQ0MjU4MDU1IiwiUGxhdGZvcm1JRCI6MiwiZXhwIjoxNjk3NzE4Mzk4LCJuYmYiOjE2ODk5NDIwOTgsImlhdCI6MTY4OTk0MjM5OH0.5d2O6yMFtyqdkkOLosYtxQoOtfsMSHn85HdQOzSX3Ok"
	ctx := testv3new.NewCtx(testcore.APIADDR, testcore.WSADDR, userID, token)
	baseCore := testcore.NewBaseCore(ctx, userID)
	ctx = mcontext.SetOperationID(ctx, utils.OperationIDGenerator())
	if err := baseCore.SendSingleMsg(ctx, recvID, 0); err != nil {
		panic(err)
	}
	if err := baseCore.SendSingleMsg(ctx, recvID, 1); err != nil {
		panic(err)
	}
}
