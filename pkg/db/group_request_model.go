package db

import (
	"errors"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertGroupRequest(groupRequest *LocalGroupRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(groupRequest).Error, "InsertGroupRequest failed")

}
func (d *DataBase) DeleteGroupRequest(groupID, userID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Where("group_id=? and user_id=?", groupID, userID).Delete(&LocalGroupRequest{}).Error, "DeleteGroupRequest failed")
}
func (d *DataBase) UpdateGroupRequest(groupRequest *LocalGroupRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Updates(groupRequest)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "_updateGroupRequest failed")
}
func (d *DataBase) GetRecvGroupApplication() ([]*LocalGroupRequest, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var groupRequestList []LocalGroupRequest
	return groupRequestList, utils.Wrap(d.conn.Where("to_user_id = ?", d.loginUserID).Find(&groupRequestList).Error, "GetRecvGroupApplication failed")
}

func (d *DataBase) GetSendGroupApplication() ([]LocalGroupRequest, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var groupRequestList []LocalGroupRequest
	return groupRequestList, utils.Wrap(d.conn.Where("user_id = ?", d.loginUserID).Find(&groupRequestList).Error, "GetSendGroupApplication failed")
}
