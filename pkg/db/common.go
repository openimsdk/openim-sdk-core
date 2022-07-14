package db

//
//import (
//	"errors"
//)
//
//func (d *DataBase) GetGroupType(groupID string) (int, error) {
//	if g, err := d.GetGroupInfoByGroupID(groupID); err == nil {
//		return int(g.GroupType), nil
//	}
//	if g, err := d.GetSuperGroupInfoByGroupID(groupID); err == nil {
//		return int(g.GroupType), nil
//	}
//	return -1, errors.New("no joined grouped " + groupID)
//}
