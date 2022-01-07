package open_im_sdk

import (
	"errors"
	"open_im_sdk/open_im_sdk/utils"
)

func (u *UserRelated) _insertFriend(friend *LocalFriend) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Create(friend).Error, "_insertFriend failed")
}

func (u *UserRelated) _deleteFriend(friendUserID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return utils.wrap(u.imdb.Where("owner_user_id=? and friend_user_id=?", u.loginUserID, friendUserID).Delete(&LocalFriend{}).Error, "_deleteFriend failed")
}

func (u *UserRelated) _updateFriend(friend *LocalFriend) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(friend)
	if t.RowsAffected == 0 {
		return utils.wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.wrap(t.Error, "_updateFriend failed")
}

func (u *UserRelated) _getAllFriendList() ([]*LocalFriend, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendList []LocalFriend
	err := utils.wrap(u.imdb.Where("owner_user_id = ?", u.loginUserID).Find(&friendList).Error,
		"_getFriendList failed")
	var transfer []*LocalFriend
	for _, v := range friendList {
		transfer = append(transfer, &v)
	}
	return transfer, err
}

func (u *UserRelated) _getFriendInfoByFriendUserID(FriendUserID string) (*LocalFriend, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friend LocalFriend
	return &friend, utils.wrap(u.imdb.Where("owner_user_id = ? AND friend_user_id = ?",
		u.loginUserID, FriendUserID).Error, "_getFriendInfoByFriendUserID failed")
}

func (u *UserRelated) _getFriendInfoList(FriendUserIDList []string) ([]LocalFriend, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendList []LocalFriend
	return friendList, utils.wrap(u.imdb.Where("friend_user_id IN ?", FriendUserIDList).Error, "_getFriendInfoListByFriendUserID failed")
}
