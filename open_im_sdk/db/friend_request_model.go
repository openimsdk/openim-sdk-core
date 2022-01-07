package db

import (
	"errors"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/utils"
)

func (u *open_im_sdk.UserRelated) _insertFriendRequest(friendRequest *open_im_sdk.LocalFriendRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Create(friendRequest).Error, "_insertFriendRequest failed")
}

func (u *open_im_sdk.UserRelated) _deleteFriendRequestBothUserID(fromUserID, toUserID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Where("from_user_id=? and to_user_id=?", fromUserID, toUserID).Delete(&open_im_sdk.LocalFriendRequest{}).Error, "deleteFriendRequest failed")
}

func (u *open_im_sdk.UserRelated) _updateFriendRequest(friendRequest *open_im_sdk.LocalFriendRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(friendRequest)
	if t.RowsAffected == 0 {
		return utils.wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.wrap(t.Error, "updateFriendRequest failed")
}

func (u *open_im_sdk.UserRelated) _getRecvFriendApplication() ([]*open_im_sdk.LocalFriendRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendRequestList []open_im_sdk.LocalFriendRequest
	err := utils.wrap(u.imdb.Where("to_user_id = ?", u.loginUserID).Find(&friendRequestList).Error, "_getLocalFriendApplication failed")

	var transfer []*open_im_sdk.LocalFriendRequest
	for _, v := range friendRequestList {
		transfer = append(transfer, &v)
	}
	return transfer, err

}

func (u *open_im_sdk.UserRelated) _getSendFriendApplication() ([]*open_im_sdk.LocalFriendRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendRequestList []open_im_sdk.LocalFriendRequest
	err := utils.wrap(u.imdb.Where("from_user_id = ?", u.loginUserID).Find(&friendRequestList).Error, "_getLocalFriendApplication failed")

	var transfer []*open_im_sdk.LocalFriendRequest
	for _, v := range friendRequestList {
		transfer = append(transfer, &v)
	}
	return transfer, err

}
