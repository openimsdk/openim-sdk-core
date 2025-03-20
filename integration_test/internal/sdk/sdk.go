package sdk

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
)

var (
	// TestSDKs SDK slice. Index is user num
	TestSDKs []*TestSDK
)

type TestSDK struct {
	UserID string
	Num    int
	SDK    *open_im_sdk.UserContext
}

func NewTestSDK(userID string, num int, loginMgr *open_im_sdk.UserContext) *TestSDK {
	return &TestSDK{
		UserID: userID,
		Num:    num,
		SDK:    loginMgr,
	}
}
