package conversation_msg

import (
	"context"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
)

type MessageController struct {
	db db_interface.DataBase
}

func (m *MessageController) BatchInsertMessageListController(ctx context.Context, MessageList []*model_struct.LocalChatLog) error {
	if len(MessageList) == 0 {
		return nil
	}
	switch MessageList[len(MessageList)-1].SessionType {
	case constant.SuperGroupChatType:
		return m.db.SuperGroupBatchInsertMessageList(ctx, MessageList, MessageList[len(MessageList)-1].RecvID)
	default:
		return m.db.BatchInsertMessageList(ctx, MessageList)
	}
}
