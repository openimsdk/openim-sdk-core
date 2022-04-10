package advanced_interface

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/sdk_struct"
)

type AdvancedFunction interface {
	MarkGroupMessageAsRead(callback open_im_sdk_callback.Base, groupID string, msgIDList, operationID string)
	DoGroupMsgReadState(groupMsgReadList []*sdk_struct.MsgStruct)
}
