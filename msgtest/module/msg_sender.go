package module

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/internal/interaction"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"sync"

	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
)

var (
	qpsCounter    int64      // 全局变量用于统计请求数
	qpsMutex      sync.Mutex // 互斥锁用于保护全局变量的并发访问
	qpsUpdateTime time.Time  // 全局变量用于记录上次更新时间
	QPSChan       chan int64 // 用于定时更新qpsCounter的channel
)

//func init() {
//	QPSChan = make(chan int64, 100)
//}

func IncrementQPS() {
	qpsMutex.Lock()
	defer qpsMutex.Unlock()

	now := time.Now()
	// 如果距离上次更新时间超过1秒，则重置计数器
	if now.Sub(qpsUpdateTime) >= time.Second {
		QPSChan <- qpsCounter
		qpsCounter = 0
		qpsUpdateTime = now
		//log.ZError(context.Background(), "QPS", nil, "qps", "timer")

	}
	qpsCounter++
}

func GetQPS() int64 {
	qpsMutex.Lock()
	defer qpsMutex.Unlock()

	return qpsCounter
}

type msgValue struct {
	SendID      string `json:"send_id"`
	RecvID      string `json:"recv_id"`
	MsgID       string `json:"msg_id"`
	OperationID string `json:"operation_id"`
	sendTime    int64  `json:"send_time"`
	Latency     int64  `json:"latency"`
}
type errorValue struct {
	err         error
	SendID      string `json:"send_id"`
	RecvID      string `json:"recv_id"`
	MsgID       string `json:"msg_id"`
	GroupID     string `json:"group_id"`
	OperationID string `json:"operation_id"`
}
type groupMessageValue struct {
	Num        int64 `json:"num"`
	LatencySum int64 `json:"latency_sum"`
	Max        int64 `json:"max"`
	Min        int64 `json:"min"`
	Latency    int64 `json:"latency"`
}

func (e *errorValue) String() string {
	return "{" + e.err.Error() + "," + e.SendID + "," + e.RecvID + "," + e.MsgID + "," + e.OperationID + "}"
}

type SendMsgUser struct {
	timeOffset              int64
	longConnMgr             *interaction.LongConnMgr
	userID                  string
	pushMsgAndMaxSeqCh      chan common.Cmd2Value
	recvPushMsgCallback     func(msg *sdkws.MsgData)
	p                       *PressureTester
	singleFailedMessageMap  map[string]*errorValue
	groupFailedMessageMap   map[string][]*errorValue
	cancelFunc              context.CancelFunc
	ctx                     context.Context
	singleSendSampleMessage map[string]*msgValue
	singleRecvSampleMessage map[string]*msgValue
	groupSendSampleNum      map[string]int
	groupRecvSampleInfo     map[string]*groupMessageValue
	groupMessage            int64
}

func (b *SendMsgUser) GetUserID() string {
	return b.userID
}

func WithRecvPushMsgCallback(callback func(msg *sdkws.MsgData)) func(core *SendMsgUser) {
	return func(core *SendMsgUser) {
		core.recvPushMsgCallback = callback
	}
}

func newIMconfig(platformID int32, apiAddr, wsAddr string) sdk_struct.IMConfig {
	return sdk_struct.IMConfig{
		PlatformID: platformID,
		ApiAddr:    apiAddr,
		WsAddr:     wsAddr,
	}
}

func newUserCtx(userID, token string, imConfig sdk_struct.IMConfig) context.Context {
	return ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID:   userID,
		Token:    token,
		IMConfig: imConfig})
}

func NewUser(userID, token string, timeOffset int64, p *PressureTester, imConfig sdk_struct.IMConfig, opts ...func(core *SendMsgUser)) *SendMsgUser {
	pushMsgAndMaxSeqCh := make(chan common.Cmd2Value, 1000)
	ctx := newUserCtx(userID, token, imConfig)
	longConnMgr := interaction.NewLongConnMgr(ctx, &ConnListner{}, nil, pushMsgAndMaxSeqCh, nil)
	core := &SendMsgUser{
		pushMsgAndMaxSeqCh:      pushMsgAndMaxSeqCh,
		longConnMgr:             longConnMgr,
		userID:                  userID,
		p:                       p,
		singleFailedMessageMap:  make(map[string]*errorValue),
		groupFailedMessageMap:   make(map[string][]*errorValue),
		singleSendSampleMessage: make(map[string]*msgValue),
		singleRecvSampleMessage: make(map[string]*msgValue),
		groupSendSampleNum:      make(map[string]int),
		groupRecvSampleInfo:     make(map[string]*groupMessageValue),
		timeOffset:              timeOffset,
		ctx:                     ctx,
	}
	for _, opt := range opts {
		opt(core)
	}
	baseCtx, cancel := context.WithCancel(ctx)
	core.cancelFunc = cancel
	go core.recvPushMsg(baseCtx)
	go core.longConnMgr.Run(baseCtx)
	return core
}

func (b *SendMsgUser) Close(ctx context.Context) {
	b.longConnMgr.Close(ctx)
	b.cancelFunc()
}

func (b *SendMsgUser) SendMsgWithContext(userID string, index int) error {
	newCtx := mcontext.SetOperationID(b.ctx, b.userID+utils.OperationIDGenerator()+userID)
	return b.SendSingleMsg(newCtx, userID, index)
}

func (b *SendMsgUser) SendGroupMsgWithContext(groupID string, index int) error {
	newCtx := mcontext.SetOperationID(b.ctx, b.userID+utils.OperationIDGenerator()+groupID)
	return b.SendGroupMsg(newCtx, groupID, index)

}

func (b *SendMsgUser) SendSingleMsg(ctx context.Context, userID string, index int) error {
	return b.sendMsg(ctx, userID, "", index, constant.SingleChatType, fmt.Sprintf("this is test msg user %s to user %s, index: %d", b.userID, userID, index))
}

func (b *SendMsgUser) BatchSendSingleMsg(ctx context.Context, userID string, index int) error {
	content := fmt.Sprintf("this is test msg user %s to user %s, index: %d", b.userID, userID, index)
	err := b.sendMsg(ctx, userID, "", index, constant.SingleChatType, content)
	if err != nil {
		log.ZError(ctx, "send msg failed", err, "userID", userID, "index", index, "content", content)
		//b.singleFailedMessageMap[content] = err
	}
	return nil
}

func (b *SendMsgUser) SendGroupMsg(ctx context.Context, groupID string, index int) error {
	return b.sendMsg(ctx, "", groupID, index, constant.SuperGroupChatType, fmt.Sprintf("this is test msg user %s to group %s, index: %d", b.userID, groupID, index))
}

func (b *SendMsgUser) BatchSendGroupMsg(ctx context.Context, groupID string, index int) error {
	content := fmt.Sprintf("this is test msg user %s to group %s, index: %d", b.userID, groupID, index)
	err := b.sendMsg(ctx, "", groupID, index, constant.SuperGroupChatType, content)
	if err != nil {
		log.ZError(ctx, "send msg failed", err, "groupID", groupID, "index", index, "content", content)
		//b.singleFailedMessageMap[content] = err
	}
	return nil
}

func (b *SendMsgUser) sendMsg(ctx context.Context, userID, groupID string, index int, sessionType int32, content string) error {
	var resp sdkws.UserSendMsgResp
	text := sdk_struct.TextElem{Content: content}
	clientMsgID := utils.GetMsgID(b.userID)
	msg := &sdkws.MsgData{
		SendID:           b.userID,
		GroupID:          groupID,
		RecvID:           userID,
		SessionType:      sessionType,
		ContentType:      constant.Text,
		SenderNickname:   b.userID,
		Content:          []byte(utils.StructToJsonString(text)),
		CreateTime:       time.Now().UnixMilli(),
		SenderPlatformID: constant.AdminPlatformID,
		ClientMsgID:      clientMsgID,
	}
	// IncrementQPS()
	now := time.Now().UnixMilli()
	if err := b.longConnMgr.SendReqWaitResp(ctx, msg, constant.SendMsg, &resp); err != nil {
		switch sessionType {
		case constant.SingleChatType:
			b.singleFailedMessageMap[clientMsgID] = &errorValue{err: err,
				SendID: b.userID, RecvID: userID, MsgID: clientMsgID, OperationID: mcontext.GetOperationID(ctx)}
			log.ZError(ctx, "send single msg failed", err, "userID", userID, "index", index, "content", content)
		case constant.SuperGroupChatType:
			b.groupFailedMessageMap[groupID] = append(b.groupFailedMessageMap[groupID], &errorValue{err: err,
				SendID: b.userID, RecvID: groupID, MsgID: clientMsgID, GroupID: groupID, OperationID: mcontext.GetOperationID(ctx)})
			log.ZError(ctx, "send group msg failed", err, "groupID", groupID, "index", index, "content", content)
		}

		return err
	}
	switch sessionType {
	case constant.SingleChatType:
		if utils.IsContain(userID, singleSampleUserList) {
			b.singleSendSampleMessage[msg.ClientMsgID] = &msgValue{
				SendID:      msg.SendID,
				RecvID:      msg.RecvID,
				MsgID:       msg.ClientMsgID,
				OperationID: mcontext.GetOperationID(ctx),
				sendTime:    msg.SendTime,
			}
		}
	case constant.SuperGroupChatType:
		b.groupSendSampleNum[groupID]++
	}

	if resp.SendTime-now > 1500 {
		log.ZWarn(ctx, "msg recv resp is too slow", nil, "sendTime", resp.SendTime, "now", now)
	}
	return nil
}

func (b *SendMsgUser) recvPushMsg(ctx context.Context) {
	for {
		select {
		case cmd := <-b.pushMsgAndMaxSeqCh:
			switch cmd.Cmd {
			case constant.CmdPushMsg:
				pushMsgs := cmd.Value.(*sdkws.PushMessages)
				for _, push := range pushMsgs.Msgs {
					for _, msg := range push.Msgs {
						if b.recvPushMsgCallback == nil {
							b.defaultRecvPushMsgCallback(cmd.Ctx, msg)
						} else {
							b.recvPushMsgCallback(msg)
						}
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (b *SendMsgUser) defaultRecvPushMsgCallback(ctx context.Context, msg *sdkws.MsgData) {
	switch msg.SessionType {
	case constant.SingleChatType:
		if utils.IsContain(msg.RecvID, singleSampleUserList) && b.userID != msg.SendID {
			b.singleRecvSampleMessage[msg.ClientMsgID] = &msgValue{
				SendID:      msg.SendID,
				RecvID:      msg.RecvID,
				MsgID:       msg.ClientMsgID,
				OperationID: mcontext.GetOperationID(ctx),
				sendTime:    msg.SendTime,
				Latency:     b.GetRelativeServerTime() - msg.SendTime,
			}
		}
	case constant.SuperGroupChatType:
		if b.userID == b.p.groupOwnerUserID[msg.GroupID] {
			b.groupMessage++
			log.ZWarn(ctx, "recv message", nil, "userID", b.userID,
				"groupOwnerID", b.p.groupOwnerUserID[msg.GroupID], "groupMessage", b.groupMessage)
			if b.groupRecvSampleInfo[msg.GroupID] == nil {
				b.groupRecvSampleInfo[msg.GroupID] = &groupMessageValue{}
			}
			latency := b.GetRelativeServerTime() - msg.SendTime
			b.groupRecvSampleInfo[msg.GroupID].Num++
			b.groupRecvSampleInfo[msg.GroupID].LatencySum += latency
			if b.groupRecvSampleInfo[msg.GroupID].Min == 0 && b.groupRecvSampleInfo[msg.GroupID].Max == 0 {
				b.groupRecvSampleInfo[msg.GroupID].Min = latency
				b.groupRecvSampleInfo[msg.GroupID].Max = latency
			}
			if latency < b.groupRecvSampleInfo[msg.GroupID].Min {
				b.groupRecvSampleInfo[msg.GroupID].Min = latency
			}
			if latency > b.groupRecvSampleInfo[msg.GroupID].Max {
				b.groupRecvSampleInfo[msg.GroupID].Max = latency
			}
		}

	}

}

func (b *SendMsgUser) GetRelativeServerTime() int64 {
	return utils.GetCurrentTimestampByMill()
}

type ConnListner struct {
}

func (c *ConnListner) OnConnecting()     {}
func (c *ConnListner) OnConnectSuccess() {}
func (c *ConnListner) OnConnectFailed(errCode int32, errMsg string) {
	// log.ZError(context.Background(), "connect failed", nil, "errCode", errCode, "errMsg", errMsg)
}
func (c *ConnListner) OnKickedOffline()    {}
func (c *ConnListner) OnUserTokenExpired() {}
