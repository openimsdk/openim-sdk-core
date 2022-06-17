package db

import (
	"errors"
	"open_im_sdk/pkg/constant"
)

func (d *DataBase) GetGroupType(groupID string) (int, error) {
	if _, err := d.GetGroupInfoByGroupID(groupID); err == nil {
		return constant.NormalGroup, nil
	}
	if _, err := d.GetSuperGroupInfoByGroupID(groupID); err == nil {
		return constant.SuperGroup, nil
	}
	return -1, errors.New("no joined grouped " + groupID)
}
