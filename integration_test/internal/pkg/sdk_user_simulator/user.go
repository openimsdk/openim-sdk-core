package sdk_user_simulator

import (
	"sync"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/tools/errs"
)

var (
	MapLock        sync.Mutex
	UserMessageMap = make(map[string]*MsgListenerCallBak)
	timeOffset     int64
)

func GetRelativeServerTime() int64 {
	return utils.GetCurrentTimestampByMill() + timeOffset
}

func InitSDK(userID string, cf sdk_struct.IMConfig) (*open_im_sdk.LoginMgr, error) {
	userForSDK := open_im_sdk.NewLoginMgr()
	var testConnListener testConnListener
	testConnListener.UserID = userID
	isInit := userForSDK.InitSDK(cf, &testConnListener)
	if !isInit {
		return nil, errs.New("sdk init failed").Wrap()
	}

	SetListener(userForSDK, userID)

	return userForSDK, nil
}

func SetListener(userForSDK *open_im_sdk.LoginMgr, userID string) {
	var testConversation conversationCallBack
	userForSDK.SetConversationListener(&testConversation)
	var testUser userCallback
	userForSDK.SetUserListener(testUser)

	msgCallBack := NewMsgListenerCallBak(userID)
	MapLock.Lock()
	UserMessageMap[userID] = msgCallBack
	MapLock.Unlock()
	userForSDK.SetAdvancedMsgListener(msgCallBack)

	var friendshipListener testFriendshipListener
	userForSDK.SetFriendshipListener(friendshipListener)

	var groupListener testGroupListener
	userForSDK.SetGroupListener(groupListener)
}
