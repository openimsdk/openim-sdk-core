package module

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/internal/interaction"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"time"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
)

type SendMsgUser struct {
	longConnMgr         *interaction.LongConnMgr
	userID              string
	pushMsgAndMaxSeqCh  chan common.Cmd2Value
	recvPushMsgCallback func(msg *sdkws.MsgData)
	failedMessageMap    map[string]error
	cancelFunc          context.CancelFunc
	ctx                 context.Context
}

func (b SendMsgUser) GetUserID() string {
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

func NewUser(userID, token string, imConfig sdk_struct.IMConfig, opts ...func(core *SendMsgUser)) *SendMsgUser {
	pushMsgAndMaxSeqCh := make(chan common.Cmd2Value, 1000)
	ctx := newUserCtx(userID, token, imConfig)
	longConnMgr := interaction.NewLongConnMgr(ctx, &ConnListner{}, nil, pushMsgAndMaxSeqCh, nil)
	core := &SendMsgUser{
		pushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh,
		longConnMgr:        longConnMgr,
		userID:             userID,
		failedMessageMap:   make(map[string]error),
		ctx:                ctx,
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
	newCtx := mcontext.SetOperationID(b.ctx, utils.OperationIDGenerator())
	return b.SendSingleMsg(newCtx, userID, index)
}

func (b *SendMsgUser) SendGroupMsgWithContext(groupID string, index int) error {
	newCtx := mcontext.SetOperationID(b.ctx, utils.OperationIDGenerator())
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
		b.failedMessageMap[content] = err
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
		b.failedMessageMap[content] = err
	}
	return nil
}

func (b *SendMsgUser) sendMsg(ctx context.Context, userID, groupID string, index int, sesstionType int32, content string) error {
	var resp sdkws.UserSendMsgResp
	text := sdk_struct.TextElem{Content: content}
	clientMsgID := utils.GetMsgID(b.userID)
	msg := &sdkws.MsgData{
		SendID:           b.userID,
		GroupID:          groupID,
		RecvID:           userID,
		SessionType:      sesstionType,
		ContentType:      constant.Text,
		SenderNickname:   b.userID,
		Content:          []byte(utils.StructToJsonString(text)),
		CreateTime:       time.Now().UnixMilli(),
		SenderPlatformID: constant.AdminPlatformID,
		ClientMsgID:      clientMsgID,
	}
	now := time.Now().UnixMilli()
	if err := b.longConnMgr.SendReqWaitResp(ctx, msg, constant.SendMsg, &resp); err != nil {
		b.failedMessageMap[clientMsgID] = err

		return err
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
							b.defaultRecvPushMsgCallback(msg)
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

func (b *SendMsgUser) defaultRecvPushMsgCallback(msg *sdkws.MsgData) {
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
