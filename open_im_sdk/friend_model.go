package open_im_sdk

import "errors"

func InsertFriendItem(friend *Friend) error {
	//return imdb.Create(friend).Error
	return wrap(imdb.Create(friend).Error, "insertFriendItem failed")
}

func delFriendItem(ownerUserID, friendUserID string) error {
	friend := Friend{OwnerUserID: ownerUserID, FriendUserID: friendUserID}
	return wrap(imdb.Delete(&friend).Error, "delFriendItem failed")
}

func updateFriendItem(friend *Friend) error {
	t := imdb.Updates(friend)
	if t.RowsAffected == 0 {
		return wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return wrap(t.Error, "updateFriendItem failed")
}
