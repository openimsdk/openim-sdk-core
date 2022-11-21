package db

import (
	"errors"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetBlackListDB() ([]*model_struct.LocalBlack, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	if d == nil {
		return nil, errors.New("database is not open")
	}
	var blackList []model_struct.LocalBlack

	err := d.conn.Find(&blackList).Error
	var transfer []*model_struct.LocalBlack
	for _, v := range blackList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}
func (d *DataBase) GetBlackListUserID() (blackListUid []string, err error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return blackListUid, utils.Wrap(d.conn.Model(&model_struct.LocalBlack{}).Select("block_user_id").Find(&blackListUid).Error, "GetBlackList failed")
}

func (d *DataBase) GetBlackInfoByBlockUserID(blockUserID string) (*model_struct.LocalBlack, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var black model_struct.LocalBlack
	return &black, utils.Wrap(d.conn.Where("owner_user_id = ? AND block_user_id = ? ",
		d.loginUserID, blockUserID).Take(&black).Error, "GetBlackInfoByBlockUserID failed")
}

func (d *DataBase) GetBlackInfoList(blockUserIDList []string) ([]*model_struct.LocalBlack, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var blackList []model_struct.LocalBlack
	if t := d.conn.Where("block_user_id IN ? ", blockUserIDList).Find(&blackList).Error; t != nil {
		return nil, utils.Wrap(t, "GetBlackInfoList failed")
	}

	var transfer []*model_struct.LocalBlack
	for _, v := range blackList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, nil
}

func (d *DataBase) InsertBlack(black *model_struct.LocalBlack) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.Create(black).Error, "InsertBlack failed")
}

func (d *DataBase) UpdateBlack(black *model_struct.LocalBlack) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	t := d.conn.Updates(black)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateBlack failed")
}

func (d *DataBase) DeleteBlack(blockUserID string) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.Where("owner_user_id=? and block_user_id=?", d.loginUserID, blockUserID).Delete(&model_struct.LocalBlack{}).Error, "DeleteBlack failed")
}
