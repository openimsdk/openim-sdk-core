package common

import (
	"github.com/mitchellh/mapstructure"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/db"
)

func GetGroupMemberListByGroupID(callback open_im_sdk_callback.Base, operationID string, db *db.DataBase, groupID string) []*db.LocalGroupMember {
	memberList, err := db.GetGroupMemberListByGroupID(groupID)
	CheckDBErrCallback(callback, err, operationID)
	return memberList
}

func MapstructureDecode(input interface{}, output interface{}, callback open_im_sdk_callback.Base, oprationID string) {
	err := mapstructure.Decode(input, output)
	CheckDataErrCallback(callback, err, oprationID)
}
