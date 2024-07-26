package sdk_user_simulator

import (
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/version"
	"github.com/openimsdk/tools/log"
)

var (
	UserMessageMap = make(map[string]*MsgListenerCallBak)
	timeOffset     int64
)

func GetRelativeServerTime() int64 {
	return utils.GetCurrentTimestampByMill() + timeOffset
}

func InitSDKAndLogin(userID, token string, cf sdk_struct.IMConfig) (*open_im_sdk.LoginMgr, error) {
	userForSDK := open_im_sdk.NewLoginMgr()
	var testConnListener testConnListener
	userForSDK.InitSDK(cf, &testConnListener)
	if err := log.InitFromConfig(userID+"_open-im-sdk-core", "", int(vars.LogLevel), true, false, cf.DataDir, 0, 24, version.Version, false); err != nil {
		return nil, err
	}
	SetListener(userForSDK, userID)

	ctx := ccontext.WithOperationID(userForSDK.BaseCtx(), utils.OperationIDGenerator())
	err := userForSDK.Login(ctx, userID, token)
	if err != nil {
		return nil, err
	}
	return userForSDK, nil
}

func SetListener(userForSDK *open_im_sdk.LoginMgr, userID string) {
	var testConversation conversationCallBack
	userForSDK.SetConversationListener(&testConversation)
	var testUser userCallback
	userForSDK.SetUserListener(testUser)

	msgCallBack := NewMsgListenerCallBak(userID)
	UserMessageMap[userID] = msgCallBack
	userForSDK.SetAdvancedMsgListener(msgCallBack)

	var friendListener testFriendListener
	userForSDK.SetFriendListener(friendListener)

	var groupListener testGroupListener
	userForSDK.SetGroupListener(groupListener)
}
