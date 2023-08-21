package sdk_user_simulator

import (
	"fmt"
	"open_im_sdk/internal/login"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

var (
	UserRouterMap map[string]*login.LoginMgr
)

var (
	TESTIP     = "59.36.173.89"
	APIADDR    = fmt.Sprintf("http://%v:10002", TESTIP)
	WSADDR     = fmt.Sprintf("ws://%v:10001", TESTIP)
	SECRET     = "openIM123"
	PLATFORMID = constant.WindowsPlatformID
	LogLevel   = uint32(5)
)

func InitSDKAndLogin(userID, token string) {
	userForSDK := login.NewLoginMgr()
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.PlatformID = int32(PLATFORMID)
	cf.WsAddr = WSADDR
	cf.DataDir = "./"
	cf.LogLevel = LogLevel
	cf.IsExternalExtensions = true
	cf.IsLogStandardOutput = true
	cf.LogFilePath = ""
	var testConnListener testConnListener
	userForSDK.InitSDK(cf, &testConnListener)
	ctx := ccontext.WithOperationID(userForSDK.BaseCtx(), utils.OperationIDGenerator())
	SetListener(userForSDK)
	userForSDK.Login(ctx, userID, token)
	UserRouterMap[userID] = userForSDK
}

func SetListener(userForSDK *login.LoginMgr) {
	var testConversation conversationCallBack
	userForSDK.SetConversationListener(&testConversation)

	var testUser userCallback
	userForSDK.SetUserListener(testUser)

	var msgCallBack MsgListenerCallBak
	userForSDK.SetAdvancedMsgListener(&msgCallBack)

	var friendListener testFriendListener
	userForSDK.SetFriendListener(friendListener)

	var groupListener testGroupListener
	userForSDK.SetGroupListener(groupListener)
}
