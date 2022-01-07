package db

import (
	"errors"
	"open_im_sdk/pkg/utils"
)

func (u *DataBase) _insertGroupRequest(groupRequest *LocalGroupRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.Wrap(u.conn.Create(groupRequest).Error, "_insertGroupRequest failed")

}
func (u *DataBase) _deleteGroupRequest(groupID, userID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.Wrap(u.conn.Where("group_id=? and user_id=?", groupID, userID).Delete(&LocalGroupRequest{}).Error, "_deleteGroupRequest failed")
}
func (u *DataBase) _updateGroupRequest(groupRequest *LocalGroupRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.conn.Updates(groupRequest)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "_updateGroupRequest failed")
}
func (u *DataBase) _getRecvGroupApplication() ([]LocalGroupRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var groupRequestList []LocalGroupRequest
	return groupRequestList, utils.Wrap(u.conn.Where("to_user_id = ?", u.loginUserID).Find(&groupRequestList).Error, "_getRecvGroupApplication failed")
}

func (u *DataBase) _getSendGroupApplication() ([]LocalGroupRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var groupRequestList []LocalGroupRequest
	return groupRequestList, utils.Wrap(u.conn.Where("user_id = ?", u.loginUserID).Find(&groupRequestList).Error, "_getSendGroupApplication failed")
}
