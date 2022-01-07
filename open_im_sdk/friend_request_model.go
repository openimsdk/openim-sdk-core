package open_im_sdk

import (
	"errors"
	"open_im_sdk/open_im_sdk/utils"
)

func (u *UserRelated) _insertFriendRequest(friendRequest *LocalFriendRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Create(friendRequest).Error, "_insertFriendRequest failed")
}

func (u *UserRelated) _deleteFriendRequestBothUserID(fromUserID, toUserID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Where("from_user_id=? and to_user_id=?", fromUserID, toUserID).Delete(&LocalFriendRequest{}).Error, "deleteFriendRequest failed")
}

func (u *UserRelated) _updateFriendRequest(friendRequest *LocalFriendRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(friendRequest)
	if t.RowsAffected == 0 {
		return utils.wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.wrap(t.Error, "updateFriendRequest failed")
}

func (u *UserRelated) _getRecvFriendApplication() ([]*LocalFriendRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendRequestList []LocalFriendRequest
	err := utils.wrap(u.imdb.Where("to_user_id = ?", u.loginUserID).Find(&friendRequestList).Error, "_getLocalFriendApplication failed")

	var transfer []*LocalFriendRequest
	for _, v := range friendRequestList {
		transfer = append(transfer, &v)
	}
	return transfer, err

}

func (u *UserRelated) _getSendFriendApplication() ([]*LocalFriendRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendRequestList []LocalFriendRequest
	err := utils.wrap(u.imdb.Where("from_user_id = ?", u.loginUserID).Find(&friendRequestList).Error, "_getLocalFriendApplication failed")

	var transfer []*LocalFriendRequest
	for _, v := range friendRequestList {
		transfer = append(transfer, &v)
	}
	return transfer, err

}
