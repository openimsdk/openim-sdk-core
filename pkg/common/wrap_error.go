package common

import (
	"github.com/mitchellh/mapstructure"
	"open_im_sdk/pkg/db"
)

func  GetGroupMemberListByGroupID(callback Base, operationID string,  db *db.DataBase, groupID string) []*db.LocalGroupMember{
	memberList, err := db.GetGroupMemberListByGroupID(groupID)
	CheckErr(callback, err, operationID)
	return memberList
}

func MapstructureDecode(input interface{}, output interface{}, callback Base, oprationID string){
	err := mapstructure.Decode(input, output)
	CheckDataErr(callback, err, oprationID)
}

