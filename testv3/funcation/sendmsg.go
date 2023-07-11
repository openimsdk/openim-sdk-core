// @Author BanTanger 2023/7/10 15:30:00
package funcation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

func SendMsg(ctx context.Context, fromID, recvID, groupID, msg string) (*sdk_struct.MsgStruct, error) {
	operationID := utils.OperationIDGenerator()
	message, err := AllLoginMgr[fromID].mgr.Conversation().CreateTextMessage(ctx, msg)
	if err != nil {
		log.Error(operationID, "CreateTextMessage ", err)
		return nil, err
	}
	o := sdkws.OfflinePushInfo{Title: "title", Desc: "desc"}
	log.Info(operationID, "SendMessage ", fromID, recvID, groupID, message.ClientMsgID)
	AllLoginMgr[fromID].mgr.Conversation().SendMessage(ctx, message, recvID, groupID, &o)
	return message, nil
}
