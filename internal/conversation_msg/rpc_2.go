package conversation_msg

import (
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/db"
)

func (c *Conversation) insertGroupMessageToLocalStorage(callback common.Base, s *db.LocalChatLog, operationID string) string {
	err := c.db.InsertMessage(s)
	common.CheckDBErr(callback, err, operationID)
	return s.ClientMsgID
}
