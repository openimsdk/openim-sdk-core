package db

import (
	_ "database/sql"
	"errors"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertFriend(friend *LocalFriend) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(friend).Error, "InsertFriend failed")
}

func (d *DataBase) DeleteFriend(friendUserID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Where("owner_user_id=? and friend_user_id=?", d.loginUserID, friendUserID).Delete(&LocalFriend{}).Error, "DeleteFriend failed")
}

func (d *DataBase) UpdateFriend(friend *LocalFriend) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	t := d.conn.Model(friend).Select("*").Updates(*friend)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")

}

func (d *DataBase) GetAllFriendList() ([]*LocalFriend, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var friendList []LocalFriend
	err := utils.Wrap(d.conn.Where("owner_user_id = ?", d.loginUserID).Find(&friendList).Error,
		"GetFriendList failed")
	var transfer []*LocalFriend
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}

func (d *DataBase) GetFriendInfoByFriendUserID(FriendUserID string) (*LocalFriend, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var friend LocalFriend
	return &friend, utils.Wrap(d.conn.Where("owner_user_id = ? AND friend_user_id = ?",
		d.loginUserID, FriendUserID).Take(&friend).Error, "GetFriendInfoByFriendUserID failed")
}

func (d *DataBase) GetFriendInfoList(friendUserIDList []string) ([]*LocalFriend, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var friendList []LocalFriend
	err := utils.Wrap(d.conn.Where("friend_user_id IN ?", friendUserIDList).Find(&friendList).Error, "GetFriendInfoListByFriendUserID failed")
	var transfer []*LocalFriend
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}
