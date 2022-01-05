package open_im_sdk

import "errors"

func (u *UserRelated) _insertGroupRequest(groupRequest *LocalGroupRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return wrap(u.imdb.Create(groupRequest).Error, "_insertGroupRequest failed")

}
func (u *UserRelated) _deleteGroupRequest(groupRequest *LocalGroupRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return Wrap(u.imdb.Delete(groupRequest).Error, "_deleteGroupRequest failed")
}
func (u *UserRelated) _updateGroupRequest(groupRequest *LocalGroupRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(groupRequest)
	if t.RowsAffected == 0 {
		return wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return wrap(t.Error, "_updateGroupRequest failed")
}
func (u *UserRelated) _getRecvGroupApplication() ([]LocalGroupRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var groupRequestList []LocalGroupRequest
	return groupRequestList, wrap(u.imdb.Where("to_user_id = ?", u.loginUserID).Find(&groupRequestList).Error, "_getRecvGroupApplication failed")
}

func (u *UserRelated) _getSendGroupApplication() ([]LocalGroupRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var groupRequestList []LocalGroupRequest
	return groupRequestList, wrap(u.imdb.Where("user_id = ?", u.loginUserID).Find(&groupRequestList).Error, "_getSendGroupApplication failed")
}
