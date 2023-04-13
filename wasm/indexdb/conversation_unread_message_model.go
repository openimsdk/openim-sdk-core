//go:build js && wasm
// +build js,wasm

package indexdb

import "context"

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

type LocalConversationUnreadMessages struct {
}

func NewLocalConversationUnreadMessages() *LocalConversationUnreadMessages {
	return &LocalConversationUnreadMessages{}
}

func (i *LocalConversationUnreadMessages) BatchInsertConversationUnreadMessageList(ctx context.Context, messageList []*model_struct.LocalConversationUnreadMessage) error {
	if messageList == nil {
		return nil
	}
	_, err := Exec(utils.StructToJsonString(messageList))
	return err
}

func (i *LocalConversationUnreadMessages) DeleteConversationUnreadMessageList(ctx context.Context, conversationID string, sendTime int64) int64 {
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
