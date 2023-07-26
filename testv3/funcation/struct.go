// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

type initLister struct {
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
