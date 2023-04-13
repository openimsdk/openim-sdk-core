//go:build !js
// +build !js

package db

import "context"

//
//import (
//	"errors"
//)
//
//func (d *DataBase) GetGroupType(ctx context.Context, groupID string) (int, error) {
//	if g, err := d.GetGroupInfoByGroupID(groupID); err == nil {
//		return int(g.GroupType), nil
//	}
//	if g, err := d.GetSuperGroupInfoByGroupID(groupID); err == nil {
//		return int(g.GroupType), nil
//	}
//	return -1, errors.New("no joined grouped " + groupID)
//}
