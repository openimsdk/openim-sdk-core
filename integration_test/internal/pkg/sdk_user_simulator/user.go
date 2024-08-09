package sdk_user_simulator

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
	"sync"
)

var (
	MapLock        sync.Mutex
	UserMessageMap = make(map[string]*MsgListenerCallBak)
	timeOffset     int64
)

func GetRelativeServerTime() int64 {
	return utils.GetCurrentTimestampByMill() + timeOffset
}

func InitSDK(ctx context.Context, userID, token string, cf sdk_struct.IMConfig) (context.Context, *open_im_sdk.LoginMgr, error) {
	userForSDK := open_im_sdk.NewLoginMgr()
	var testConnListener testConnListener
	isInit := userForSDK.InitSDK(cf, &testConnListener)
	if !isInit {
		return nil, nil, errs.New("sdk init failed").Wrap()
	}

	SetListener(userForSDK, userID)

	ctx = userForSDK.Context()
	ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
	ctx = mcontext.SetOpUserID(ctx, userID)
	if err := userForSDK.InitMgr(ctx, userID, token); err != nil {
		return nil, nil, err
	}

	return ctx, userForSDK, nil
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

	var friendListener testFriendListener
	userForSDK.SetFriendListener(friendListener)

	var groupListener testGroupListener
	userForSDK.SetGroupListener(groupListener)
}
