package open_im_sdk

import "errors"

func (u *UserRelated) _insertFriendRequestData(friendRequest *FriendRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	return wrap(u.imdb.Create(friendRequest).Error, "insertFriendRequestData failed")
}

func (u *UserRelated) _deleteFriendRequestBothUserID(fromUserID, toUserID string) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	friendRequest := FriendRequest{FromUserID: fromUserID, ToUserID: toUserID}
	return wrap(u.imdb.Delete(&friendRequest).Error, "deleteFriendRequest failed")
}

func (u *UserRelated) _updateFriendRequest(friendRequest *FriendRequest) error {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	t := u.imdb.Updates(friendRequest)
	if t.RowsAffected == 0 {
		return wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return wrap(t.Error, "updateFriendRequest failed")
}

func (u *UserRelated) _getRecvFriendApplication() ([]FriendRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendRequestList []FriendRequest
	return friendRequestList, wrap(u.imdb.Where("to_user_id = ?", u.loginUserID).Find(&friendRequestList).Error, "_getLocalFriendApplication failed")
}

func (u *UserRelated) _getSendFriendApplication() ([]FriendRequest, error) {
	u.mRWMutex.Lock()
	defer u.mRWMutex.Unlock()
	var friendRequestList []FriendRequest
	return friendRequestList, wrap(u.imdb.Where("from_user_id = ?", u.loginUserID).Find(&friendRequestList).Error, "_getLocalFriendApplication failed")
}
