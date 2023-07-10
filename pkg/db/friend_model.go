package db

import (
	_ "database/sql"
	"errors"
	"fmt"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertFriend(friend *model_struct.LocalFriend) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.Create(friend).Error, "InsertFriend failed")
}

func (d *DataBase) DeleteFriendDB(friendUserID string) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.Where("owner_user_id=? and friend_user_id=?", d.loginUserID, friendUserID).Delete(&model_struct.LocalFriend{}).Error, "DeleteFriend failed")
}

func (d *DataBase) UpdateFriend(friend *model_struct.LocalFriend) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()

	t := d.conn.Model(friend).Select("*").Updates(*friend)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")

}

func (d *DataBase) GetAllFriendList() ([]*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendList []model_struct.LocalFriend
	err := utils.Wrap(d.conn.Where("owner_user_id = ?", d.loginUserID).Find(&friendList).Error,
		"GetFriendList failed")
	var transfer []*model_struct.LocalFriend
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}
func (d *DataBase) SearchFriendList(keyword string, isSearchUserID, isSearchNickname, isSearchRemark bool) ([]*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var count int
	var friendList []model_struct.LocalFriend
	var condition string
	if isSearchUserID {
		condition = fmt.Sprintf("friend_user_id like %q ", "%"+keyword+"%")
		count++
	}
	if isSearchNickname {
		if count > 0 {
			condition += "or "
		}
		condition += fmt.Sprintf("name like %q ", "%"+keyword+"%")
		count++
	}
	if isSearchRemark {
		if count > 0 {
			condition += "or "
		}
		condition += fmt.Sprintf("remark like %q ", "%"+keyword+"%")
	}
	err := d.conn.Where(condition).Order("create_time DESC").Find(&friendList).Error
	var transfer []*model_struct.LocalFriend
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "SearchFriendList failed ")

}

func (d *DataBase) GetFriendInfoByFriendUserID(FriendUserID string) (*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friend model_struct.LocalFriend
	return &friend, utils.Wrap(d.conn.Where("owner_user_id = ? AND friend_user_id = ?",
		d.loginUserID, FriendUserID).Take(&friend).Error, "GetFriendInfoByFriendUserID failed")
}

func (d *DataBase) GetFriendInfoList(friendUserIDList []string) ([]*model_struct.LocalFriend, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendList []model_struct.LocalFriend
	err := utils.Wrap(d.conn.Where("friend_user_id IN ?", friendUserIDList).Find(&friendList).Error, "GetFriendInfoListByFriendUserID failed")
	var transfer []*model_struct.LocalFriend
	for _, v := range friendList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}
