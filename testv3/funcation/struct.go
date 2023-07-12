// @Author BanTanger 2023/7/10 16:01
package funcation

import (
	"open_im_sdk/internal/login"
	"time"
)

type Users struct {
	Uid      string
	Nickname string
	FaceUrl  string
}

type BaseSuccessFailed struct {
	successData string
	errCode     int
	errMsg      string
	funcName    string
	time        time.Time
}

type testInitLister struct {
}

type conversationCallBack struct {
	SyncFlag int
}

type userCallback struct {
}

type MsgListenerCallBak struct {
}

type testFriendListener struct {
	x int
}

type testGroupListener struct {
}

type CoreNode struct {
	Token             string
	UserID            string
	Mgr               *login.LoginMgr
	sendMsgSuccessNum uint32
	sendMsgFailedNum  uint32
	idx               int
}

type TestSendMsgCallBack struct {
	msg         string
	OperationID string
	sendID      string
	recvID      string
	msgID       string
	sendTime    int64
	recvTime    int64
	groupID     string
}

type SendRecvTime struct {
	SendTime             int64
	SendSeccCallbackTime int64
	RecvTime             int64
	SendIDRecvID         string
}
