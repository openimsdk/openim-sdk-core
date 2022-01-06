package open_im_sdk

import "errors"

func (u *UserRelated) _insertFriend(friend *LocalFriend) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return wrap(u.imdb.Create(friend).Error, "_insertFriend failed")
}

func (u *UserRelated) _deleteFriend(friendUserID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	//friend := LocalFriend{OwnerUserID: u.loginUserID, FriendUserID: friendUserID}
	err := u.imdb.Model(&LocalFriend{}).Where("owner_user_id=? and friend_user_id=?", u.loginUserID, friendUserID).Delete(&LocalFriend{}).Error
	return wrap(err, "_deleteFriend failed")
}

func (u *UserRelated) _updateFriend(friend *LocalFriend) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(friend)
	if t.RowsAffected == 0 {
		return wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return wrap(t.Error, "_updateFriend failed")
}

func (u *UserRelated) _getAllFriendList() ([]*LocalFriend, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendList []LocalFriend
	err := wrap(u.imdb.Where("owner_user_id = ?", u.loginUserID).Find(&friendList).Error,
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
	return &friend, wrap(u.imdb.Where("owner_user_id = ? AND friend_user_id = ?",
		u.loginUserID, FriendUserID).Error, "_getFriendInfoByFriendUserID failed")
}

func (u *UserRelated) _getFriendInfoList(FriendUserIDList []string) ([]LocalFriend, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendList []LocalFriend
	return friendList, wrap(u.imdb.Where("friend_user_id IN ?", FriendUserIDList).Error, "_getFriendInfoListByFriendUserID failed")
}
