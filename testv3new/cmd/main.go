package main

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/testv3new"
	"open_im_sdk/testv3new/testcore"
)

func main() {
	if err := log.InitFromConfig("sdk.log", "sdk", 3, true, false, "", 2); err != nil {
		panic(err)
	}
	userID := "4844258055"
	recvID := "4950983283"
	manager := testv3new.NewRegisterManager()
	token, _ := manager.GetToken(userID)
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
