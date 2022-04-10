package common

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/db"

	"github.com/mitchellh/mapstructure"
)

// 带有错误码的error
type CodeError struct {
	Code int
	Msg  string
}

func (e *CodeError) Error() string {
	if e == nil {
		return ""
	}
	return e.Msg
}

// 使用特定的code和消息构造一个CodeError
func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

func GetGroupMemberListByGroupID(callback open_im_sdk_callback.Base, operationID string, db *db.DataBase, groupID string) []*db.LocalGroupMember {
	memberList, err := db.GetGroupMemberListByGroupID(groupID)
	CheckDBErrCallback(callback, err, operationID)
	return memberList
}

func MapstructureDecode(input interface{}, output interface{}, callback open_im_sdk_callback.Base, oprationID string) {
	err := mapstructure.Decode(input, output)
	CheckDataErrCallback(callback, err, oprationID)
}
