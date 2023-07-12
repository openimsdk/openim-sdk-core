// @Author BanTanger 2023/7/10 15:30:00
package funcation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"

	"open_im_sdk/sdk_struct"
)

func SendMsg(ctx context.Context, fromID, recvID, groupID, msg string) (*sdk_struct.MsgStruct, error) {
	message, err := AllLoginMgr[fromID].Mgr.Conversation().CreateTextMessage(ctx, msg)
	if err != nil {
		log.ZError(ctx, "CreateTextMessage ", err)
		return nil, err
	}
	o := sdkws.OfflinePushInfo{Title: "title", Desc: "desc"}
	log.ZInfo(ctx, "SendMessage ", "fromID", fromID, "recvID", recvID, "groupID", groupID, "message", message)
	AllLoginMgr[fromID].Mgr.Conversation().SendMessage(ctx, message, recvID, groupID, &o)
	return message, nil
}
