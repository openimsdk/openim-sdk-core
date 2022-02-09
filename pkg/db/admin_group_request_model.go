package db

import (
	"errors"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertAdminGroupRequest(groupRequest *LocalAdminGroupRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(groupRequest).Error, "InsertAdminGroupRequest failed")
}

func (d *DataBase) DeleteAdminGroupRequest(groupID, userID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Where("group_id=? and user_id=?", groupID, userID).Delete(&LocalAdminGroupRequest{}).Error, "DeleteAdminGroupRequest failed")
}

func (d *DataBase) UpdateAdminGroupRequest(groupRequest *LocalAdminGroupRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	t := d.conn.Model(groupRequest).Select("*").Updates(*groupRequest)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")
}

func (d *DataBase) GetAdminGroupApplication() ([]*LocalAdminGroupRequest, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var groupRequestList []LocalAdminGroupRequest
	err := utils.Wrap(d.conn.Find(&groupRequestList).Error, "")
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var transfer []*LocalAdminGroupRequest
	for _, v := range groupRequestList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, nil
}
