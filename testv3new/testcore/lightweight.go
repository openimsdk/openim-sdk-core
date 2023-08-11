package testcore

import (
	"context"
	"fmt"
	"open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"

	"time"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
)

type BaseCore struct {
	longConnMgr         *interaction.LongConnMgr
	userID              string
	platformID          int32
	pushMsgAndMaxSeqCh  chan common.Cmd2Value
	recvPushMsgCallback func(msg *sdkws.MsgData)
	failedMessageMap    map[string]error
}

func (b BaseCore) GetUserID() string {
	return b.userID
}

func WithRecvPushMsgCallback(callback func(msg *sdkws.MsgData)) func(core *BaseCore) {
	return func(core *BaseCore) {
		core.recvPushMsgCallback = callback
	}
}

func NewBaseCore(ctx context.Context, userID string, platformID int32, opts ...func(core *BaseCore)) *BaseCore {
	pushMsgAndMaxSeqCh := make(chan common.Cmd2Value, 1000)
	longConnMgr := interaction.NewLongConnMgr(ctx, &ConnListner{}, nil, pushMsgAndMaxSeqCh, nil)
	core := &BaseCore{
		pushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh,
		longConnMgr:        longConnMgr,
		userID:             userID,
		platformID:         platformID,
		failedMessageMap:   make(map[string]error),
	}
	for _, opt := range opts {
		opt(core)
	}
	go core.recvPushMsg()
	go core.longConnMgr.Run(ctx)
	return core
}

func (b *BaseCore) Close(ctx context.Context) {
	b.longConnMgr.Close(ctx)
}

func (b *BaseCore) SendSingleMsg(ctx context.Context, userID string, index int) error {
	return b.sendMsg(ctx, userID, "", index, constant.SingleChatType, fmt.Sprintf("this is test msg user %s to user %s, index: %d", b.userID, userID, index))
}
func (b *BaseCore) BatchSendSingleMsg(ctx context.Context, userID string, index int) error {
	content := fmt.Sprintf("this is test msg user %s to user %s, index: %d", b.userID, userID, index)
	err := b.sendMsg(ctx, userID, "", index, constant.SingleChatType, content)
	if err != nil {
		log.ZError(ctx, "send msg failed", err, "userID", userID, "index", index, "content", content)
		b.failedMessageMap[content] = err
	}
	return nil
}

func (b *BaseCore) SendGroupMsg(ctx context.Context, groupID string, index int) error {
	return b.sendMsg(ctx, "", groupID, index, constant.SuperGroupChatType, fmt.Sprintf("this is test msg user %s to group %s, index: %d", b.userID, groupID, index))
}
func (b *BaseCore) BatchSendGroupMsg(ctx context.Context, groupID string, index int) error {
	content := fmt.Sprintf("this is test msg user %s to group %s, index: %d", b.userID, groupID, index)
	err := b.sendMsg(ctx, "", groupID, index, constant.SuperGroupChatType, content)
	if err != nil {
		log.ZError(ctx, "send msg failed", err, "groupID", groupID, "index", index, "content", content)
		b.failedMessageMap[content] = err
	}
	return nil
}

func (b *BaseCore) sendMsg(ctx context.Context, userID, groupID string, index int, sesstionType int32, content string) error {
	var resp sdkws.UserSendMsgResp
	text := sdk_struct.TextElem{Content: content}
	msg := &sdkws.MsgData{
		SendID:           b.userID,
		GroupID:          groupID,
		RecvID:           userID,
		SessionType:      sesstionType,
		ContentType:      constant.Text,
		SenderNickname:   b.userID,
		Content:          []byte(utils.StructToJsonString(text)),
		CreateTime:       time.Now().UnixMilli(),
		SenderPlatformID: b.platformID,
		ClientMsgID:      utils.GetMsgID(b.userID),
	}
	now := time.Now().UnixMilli()
	if err := b.longConnMgr.SendReqWaitResp(ctx, msg, constant.SendMsg, &resp); err != nil {
		return err
	}
	if resp.SendTime-now > 1500 {
		log.ZWarn(ctx, "msg recv resp is too slow", nil, "sendTime", resp.SendTime, "now", now)
	}
	return nil
}

func (b *BaseCore) recvPushMsg() {
	for {
		cmd := <-b.pushMsgAndMaxSeqCh
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
	}
}

func (b *BaseCore) defaultRecvPushMsgCallback(msg *sdkws.MsgData) {
}
