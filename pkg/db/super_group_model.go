package db

import "open_im_sdk/pkg/utils"

func (d *DataBase) GetJoinedSuperGroupList() ([]*LocalGroup, error) {
	return nil, nil
}

func (d *DataBase) GetJoinedSuperGroupIDList() ([]string, error) {
	return nil, nil
}

func (d *DataBase) InsertSuperGroup(groupInfo *LocalGroup) error {
	return utils.Wrap(d.conn.Table("super_groups").Create(groupInfo).Error, "InsertSuperGroup failed")
}
