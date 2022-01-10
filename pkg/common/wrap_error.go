package common

import (
	"open_im_sdk/pkg/db"
)

func  GetGroupMemberListByGroupID(callback Base, OperationID string,  db *db.DataBase, groupID string) []*db.LocalGroupMember{
	memberList, err := db.GetGroupMemberListByGroupID(groupID)
	CheckErr(callback, err, OperationID)
	return memberList
}
