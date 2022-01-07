package db

import (
	"errors"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/utils"
)

func (u *open_im_sdk.UserRelated) _insertGroupRequest(groupRequest *open_im_sdk.LocalGroupRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Create(groupRequest).Error, "_insertGroupRequest failed")

}
func (u *open_im_sdk.UserRelated) _deleteGroupRequest(groupID, userID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.Wrap(u.imdb.Where("group_id=? and user_id=?", groupID, userID).Delete(&open_im_sdk.LocalGroupRequest{}).Error, "_deleteGroupRequest failed")
}
func (u *open_im_sdk.UserRelated) _updateGroupRequest(groupRequest *open_im_sdk.LocalGroupRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(groupRequest)
	if t.RowsAffected == 0 {
		return utils.wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.wrap(t.Error, "_updateGroupRequest failed")
}
func (u *open_im_sdk.UserRelated) _getRecvGroupApplication() ([]open_im_sdk.LocalGroupRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var groupRequestList []open_im_sdk.LocalGroupRequest
	return groupRequestList, utils.wrap(u.imdb.Where("to_user_id = ?", u.loginUserID).Find(&groupRequestList).Error, "_getRecvGroupApplication failed")
}

func (u *open_im_sdk.UserRelated) _getSendGroupApplication() ([]open_im_sdk.LocalGroupRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var groupRequestList []open_im_sdk.LocalGroupRequest
	return groupRequestList, utils.wrap(u.imdb.Where("user_id = ?", u.loginUserID).Find(&groupRequestList).Error, "_getSendGroupApplication failed")
}
