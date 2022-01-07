package open_im_sdk

import (
	"errors"
	"open_im_sdk/open_im_sdk/utils"
)

func (u *UserRelated) _insertGroupRequest(groupRequest *LocalGroupRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Create(groupRequest).Error, "_insertGroupRequest failed")

}
func (u *UserRelated) _deleteGroupRequest(groupID, userID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.Wrap(u.imdb.Where("group_id=? and user_id=?", groupID, userID).Delete(&LocalGroupRequest{}).Error, "_deleteGroupRequest failed")
}
func (u *UserRelated) _updateGroupRequest(groupRequest *LocalGroupRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(groupRequest)
	if t.RowsAffected == 0 {
		return utils.wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.wrap(t.Error, "_updateGroupRequest failed")
}
func (u *UserRelated) _getRecvGroupApplication() ([]LocalGroupRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var groupRequestList []LocalGroupRequest
	return groupRequestList, utils.wrap(u.imdb.Where("to_user_id = ?", u.loginUserID).Find(&groupRequestList).Error, "_getRecvGroupApplication failed")
}

func (u *UserRelated) _getSendGroupApplication() ([]LocalGroupRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var groupRequestList []LocalGroupRequest
	return groupRequestList, utils.wrap(u.imdb.Where("user_id = ?", u.loginUserID).Find(&groupRequestList).Error, "_getSendGroupApplication failed")
}
