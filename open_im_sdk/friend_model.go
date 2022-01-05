package open_im_sdk

import "errors"

func (u *UserRelated) _insertFriend(friend *Friend) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return wrap(u.imdb.Create(friend).Error, "_insertFriend failed")
}

func (u *UserRelated) _deleteFriend(friendUserID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	friend := Friend{OwnerUserID: u.loginUserID, FriendUserID: friendUserID}
	return wrap(u.imdb.Delete(&friend).Error, "_deleteFriend failed")
}

func (u *UserRelated) _updateFriend(friend *Friend) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(friend)
	if t.RowsAffected == 0 {
		return wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return wrap(t.Error, "_updateFriend failed")
}

func (u *UserRelated) _getFriendList() ([]Friend, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendList []Friend
	return friendList, wrap(u.imdb.Where("owner_user_id = ?", u.loginUserID).Find(&friendList).Error,
		"_getFriendList failed")
}

func (u *UserRelated) _getFriendInfoByFriendUserID(FriendUserID string) (*Friend, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friend Friend
	return &friend, wrap(u.imdb.Where("owner_user_id = ? AND friend_user_id = ?",
		u.loginUserID, FriendUserID).Error, "_getFriendInfoByFriendUserID failed")
}
