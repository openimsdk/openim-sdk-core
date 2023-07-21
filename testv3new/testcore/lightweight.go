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
	longConnMgr        *interaction.LongConnMgr
	userID             string
	recvMsgNum         int
	platformID         int32
	pushMsgAndMaxSeqCh chan common.Cmd2Value
}

func NewBaseCore(userID string) *BaseCore {
	ctx := context.Background()
	pushMsgAndMaxSeqCh := make(chan common.Cmd2Value, 1000)
	longConnMgr := interaction.NewLongConnMgr(ctx, &ConnListner{}, nil, pushMsgAndMaxSeqCh, nil)
	return &BaseCore{
		pushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh,
		longConnMgr:        longConnMgr,
		userID:             userID,
		platformID:         constant.AndroidPlatformID,
	}
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
	if err := b.longConnMgr.SendReqWaitResp(ctx, msg, constant.SendMsg, &resp); err != nil {
		return err
	}
	if resp.SendTime-now > 1500 {
		log.ZWarn(ctx, "msg recv resp is too slow", nil, "sendTime", resp.SendTime, "now", now)
	}
	return nil
}

func (b *BaseCore) recvPushMsg() {
	// for {
	// 	cmd := <-b.PushMsgAndMaxSeqCh

	// }
}

func (b *BaseCore) recvMsgCallback(msg *sdkws.MsgData) {
	b.recvMsgNum++
}

func (b *BaseCore) GetRecvMsgNum() int {
	return b.recvMsgNum
}
