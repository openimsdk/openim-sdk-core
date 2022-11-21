package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

type LocalConversationUnreadMessages struct {
}

func (i *LocalConversationUnreadMessages) BatchInsertConversationUnreadMessageList(messageList []*model_struct.LocalConversationUnreadMessage) error {
	if messageList == nil {
		return nil
	}
	_, err := Exec(utils.StructToJsonString(messageList))
	return err
}

func (i *LocalConversationUnreadMessages) DeleteConversationUnreadMessageList(conversationID string, sendTime int64) int64 {
	deleteRows, err := Exec(conversationID, sendTime)
	if err != nil {
		return 0
	} else {
		if v, ok := deleteRows.(float64); ok {
			var result int64
			result = int64(v)
			return result
		} else {
			return 0
		}
	}
}
