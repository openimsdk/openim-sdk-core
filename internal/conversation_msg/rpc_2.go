package conversation_msg

import (
	"encoding/json"
	"open_im_sdk/pkg/common"
)

func (c *Conversation) insertGroupMessageToLocalStorage(callback common.Base, message, groupID, sender string, operationID string) string {

	if err := c.db.InsertMessage()(&s); err != nil {
		callback.OnError(201, err.Error())
	} else {
		callback.OnSuccess("")
	}

	return s.ClientMsgID
}
