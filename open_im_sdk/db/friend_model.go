package db

import (
	"errors"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/utils"
)

func (u *open_im_sdk.UserRelated) _insertFriend(friend *open_im_sdk.LocalFriend) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Create(friend).Error, "_insertFriend failed")
}

func (u *open_im_sdk.UserRelated) _deleteFriend(friendUserID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Where("owner_user_id=? and friend_user_id=?", u.loginUserID, friendUserID).Delete(&open_im_sdk.LocalFriend{}).Error, "_deleteFriend failed")
}

func (u *open_im_sdk.UserRelated) _updateFriend(friend *open_im_sdk.LocalFriend) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(friend)
	if t.RowsAffected == 0 {
		return utils.wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.wrap(t.Error, "_updateFriend failed")
}

func (u *open_im_sdk.UserRelated) _getAllFriendList() ([]*open_im_sdk.LocalFriend, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendList []open_im_sdk.LocalFriend
	err := utils.wrap(u.imdb.Where("owner_user_id = ?", u.loginUserID).Find(&friendList).Error,
		"_getFriendList failed")
	var transfer []*open_im_sdk.LocalFriend
	for _, v := range friendList {
		transfer = append(transfer, &v)
	}
	return transfer, err
}

func (u *open_im_sdk.UserRelated) _getFriendInfoByFriendUserID(FriendUserID string) (*open_im_sdk.LocalFriend, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friend open_im_sdk.LocalFriend
	return &friend, utils.wrap(u.imdb.Where("owner_user_id = ? AND friend_user_id = ?",
		u.loginUserID, FriendUserID).Error, "_getFriendInfoByFriendUserID failed")
}

func (u *open_im_sdk.UserRelated) _getFriendInfoList(FriendUserIDList []string) ([]open_im_sdk.LocalFriend, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendList []open_im_sdk.LocalFriend
	return friendList, utils.wrap(u.imdb.Where("friend_user_id IN ?", FriendUserIDList).Error, "_getFriendInfoListByFriendUserID failed")
}
