package common

import (
	"github.com/mitchellh/mapstructure"
	"open_im_sdk/pkg/db"
)

func GetGroupMemberListByGroupID(callback Base, operationID string, db *db.DataBase, groupID string) []*db.LocalGroupMember {
	memberList, err := db.GetGroupMemberListByGroupID(groupID)
	CheckDBErrCallback(callback, err, operationID)
	return memberList
}

func MapstructureDecode(input interface{}, output interface{}, callback Base, oprationID string) {
	err := mapstructure.Decode(input, output)
	CheckDataErrCallback(callback, err, oprationID)
}
