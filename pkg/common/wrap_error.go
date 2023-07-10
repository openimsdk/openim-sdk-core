package common

//import (
//	"github.com/mitchellh/mapstructure"
//	"open_im_sdk/open_im_sdk_callback"
//	"open_im_sdk/pkg/db"
//	"open_im_sdk/pkg/db/model_struct"
//)
//
//func GetGroupMemberListByGroupID(callback open_im_sdk_callback.Base, operationID string, db *db.DataBase, groupID string) []*model_struct.LocalGroupMember {
//	memberList, err := db.GetGroupMemberListByGroupID(groupID)
//	CheckDBErrCallback(callback, err, operationID)
//	return memberList
//}
//
//func MapstructureDecode(input interface{}, output interface{}, callback open_im_sdk_callback.Base, oprationID string) {
//	err := mapstructure.Decode(input, output)
//	CheckDataErrCallback(callback, err, oprationID)
//}
