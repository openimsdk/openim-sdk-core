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

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

type BaseCore struct {
	longConnMgr         *interaction.LongConnMgr
	userID              string
	platformID          int32
	pushMsgAndMaxSeqCh  chan common.Cmd2Value
	recvPushMsgCallback func(msg *sdkws.MsgData)
	wsUrl               string
	recvMap             map[string]int
}

func (b BaseCore) GetRecvMap() map[string]int {
	if b.recvMap != nil {
		return b.recvMap
	}
	return nil
}

func WithRecvPushMsgCallback(callback func(msg *sdkws.MsgData)) func(core *BaseCore) {
	return func(core *BaseCore) {
		core.recvPushMsgCallback = callback
	}
}

func NewBaseCore(ctx context.Context, userID string, opts ...func(core *BaseCore)) *BaseCore {
	pushMsgAndMaxSeqCh := make(chan common.Cmd2Value, 1000)
	longConnMgr := interaction.NewLongConnMgr(ctx, &ConnListner{}, nil, pushMsgAndMaxSeqCh, nil)
	core := &BaseCore{
		pushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh,
		longConnMgr:        longConnMgr,
		userID:             userID,
		platformID:         constant.AndroidPlatformID,
		recvMap:            map[string]int{},
	}
	for _, opt := range opts {
		opt(core)
	}
	go core.recvPushMsg()
	go core.longConnMgr.Run(ctx)
	return core
}

func (b *BaseCore) SendSingleMsg(ctx context.Context, userID string, index int) error {
	return b.sendMsg(ctx, userID, "", index)
}

func (b *BaseCore) SendGroupMsg(ctx context.Context, groupID string, index int) error {
	return b.sendMsg(ctx, "", groupID, index)
}

func (b *BaseCore) sendMsg(ctx context.Context, userID, groupID string, index int) error {
	var resp sdkws.UserSendMsgResp
	var sesstionType int32
	var content string
	if userID != "" {
		sesstionType = constant.SingleChatType
		content = fmt.Sprintf("this is test msg user %s to user %s, index: %d", b.userID, userID, index)
	} else {
		sesstionType = constant.SuperGroupChatType
		content = fmt.Sprintf("this is test msg user %s to group %s, index: %d", b.userID, groupID, index)
	}
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
	// time.Sleep(60 * time.Second)
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
	if b.userID == msg.RecvID {
		b.recvMap[msg.SendID+"_"+msg.RecvID]++
	}
}
