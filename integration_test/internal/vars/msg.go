package vars

import (
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"sync/atomic"
	"time"
)

type StatMsg struct {
	CostTime    time.Duration
	ReceiveTime time.Time
	Msg         *sdk_struct.MsgStruct
}

var (
	SendMsgCount     atomic.Int64
	RecvMsgConsuming chan *StatMsg
)
