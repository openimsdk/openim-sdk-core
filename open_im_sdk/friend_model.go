package open_im_sdk

import "errors"

func (u *UserRelated) InsertFriendItem(friend *Friend) error {
	return wrap(u.imdb.Create(friend).Error, "insertFriendItem failed")
}

func (u *UserRelated) delFriendItem(ownerUserID, friendUserID string) error {
	friend := Friend{OwnerUserID: ownerUserID, FriendUserID: friendUserID}
	return wrap(u.imdb.Delete(&friend).Error, "delFriendItem failed")
}

func (u *UserRelated) updateFriendItem(friend *Friend) error {
	t := u.imdb.Updates(friend)
	if t.RowsAffected == 0 {
		return wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return wrap(t.Error, "updateFriendItem failed")
}

func (u *UserRelated) getFriendList() ([]Friend, error) {
	var friendList []Friend
	return friendList, wrap(u.imdb.Where("owner_user_id = ?", u.LoginUid).Find(&friendList).Error, "getFriendListByUserID failed")
}
