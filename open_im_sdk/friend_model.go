package open_im_sdk

import "errors"

func (u *UserRelated) InsertFriendItem(friend *Friend) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return wrap(u.imdb.Create(friend).Error, "insertFriendItem failed")
}

func (u *UserRelated) delFriendItem(friendUserID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	friend := Friend{OwnerUserID: u.loginUserID, FriendUserID: friendUserID}
	return wrap(u.imdb.Delete(&friend).Error, "delFriendItem failed")
}

func (u *UserRelated) updateFriendItem(friend *Friend) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(friend)
	if t.RowsAffected == 0 {
		return wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return wrap(t.Error, "updateFriendItem failed")
}

func (u *UserRelated) getFriendList() ([]Friend, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendList []Friend
	return friendList, wrap(u.imdb.Where("owner_user_id = ?", u.loginUserID).Find(&friendList).Error,
		"getFriendList failed")
}

func (u *UserRelated) getFriendInfoByFriendUserID(FriendUserID string) (*Friend, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friend Friend
	return &friend, wrap(u.imdb.Where("owner_user_id = ? AND friend_user_id = ?",
		u.loginUserID, FriendUserID).Error, "getFriendInfoByFriendUserID failed")
}
