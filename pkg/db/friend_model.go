package db

import (
	_ "database/sql"
	"errors"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) _insertFriend(friend *LocalFriend) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(friend).Error, "_insertFriend failed")
}

func (d *DataBase) _deleteFriend(friendUserID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Where("owner_user_id=? and friend_user_id=?", d.loginUserID, friendUserID).Delete(&LocalFriend{}).Error, "_deleteFriend failed")
}

func (d *DataBase) _updateFriend(friend *LocalFriend) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Updates(friend)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "_updateFriend failed")
}

func (d *DataBase) _getAllFriendList() ([]*LocalFriend, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var friendList []LocalFriend
	err := utils.Wrap(d.conn.Where("owner_user_id = ?", d.loginUserID).Find(&friendList).Error,
		"_getFriendList failed")
	var transfer []*LocalFriend
	for _, v := range friendList {
		transfer = append(transfer, &v)
	}
	return transfer, err
}

func (d *DataBase) _getFriendInfoByFriendUserID(FriendUserID string) (*LocalFriend, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var friend LocalFriend
	return &friend, utils.Wrap(d.conn.Where("owner_user_id = ? AND friend_user_id = ?",
		d.loginUserID, FriendUserID).Error, "_getFriendInfoByFriendUserID failed")
}

func (d *DataBase) _getFriendInfoList(FriendUserIDList []string) ([]LocalFriend, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var friendList []LocalFriend
	return friendList, utils.Wrap(d.conn.Where("friend_user_id IN ?", FriendUserIDList).Error, "_getFriendInfoListByFriendUserID failed")
}
