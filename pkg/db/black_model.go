package db

import (
	"errors"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetBlackList() ([]*LocalBlack, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var blackList []LocalBlack

	err := d.conn.Find(&blackList).Error
	var transfer []*LocalBlack
	for _, v := range blackList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err

}
func (d *DataBase) GetBlackListUserID() (blackListUid []string, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	return blackListUid, utils.Wrap(d.conn.Model(&LocalBlack{}).Select("block_user_id").Find(&blackListUid).Error, "GetBlackList failed")
}

func (d *DataBase) GetBlackInfoByBlockUserID(blockUserID string) (*LocalBlack, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var black LocalBlack
	return &black, utils.Wrap(d.conn.Where("owner_user_id = ? AND block_user_id = ? ",
		d.loginUserID, blockUserID).Scan(&black).Error, "GetBlackInfoByBlockUserID failed")
}

func (d *DataBase) GetBlackInfoList(blockUserIDList []string) ([]LocalBlack, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var black []LocalBlack
	return black, utils.Wrap(d.conn.Where("block_user_id IN ? ", blockUserIDList).Find(&black).Error, "GetBlackInfoList failed")
}

func (d *DataBase) InsertBlack(black *LocalBlack) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(black).Error, "InsertBlack failed")
}

func (d *DataBase) UpdateBlack(black *LocalBlack) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Updates(black)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateBlack failed")
}

func (d *DataBase) DeleteBlack(blockUserID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Where("owner_user_id=? and block_user_id=?", d.loginUserID, blockUserID).Delete(&LocalBlack{}).Error, "DeleteBlack failed")
}
