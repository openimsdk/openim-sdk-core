package db

import (
	"errors"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertFriendRequest(friendRequest *model_struct.LocalFriendRequest) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.Create(friendRequest).Error, "InsertFriendRequest failed")
}

func (d *DataBase) DeleteFriendRequestBothUserID(fromUserID, toUserID string) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	return utils.Wrap(d.conn.Where("from_user_id=? and to_user_id=?", fromUserID, toUserID).Delete(&model_struct.LocalFriendRequest{}).Error, "DeleteFriendRequestBothUserID failed")
}

func (d *DataBase) UpdateFriendRequest(friendRequest *model_struct.LocalFriendRequest) error {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	t := d.conn.Model(friendRequest).Select("*").Updates(*friendRequest)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")
}

func (d *DataBase) GetRecvFriendApplication() ([]*model_struct.LocalFriendRequest, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendRequestList []model_struct.LocalFriendRequest
	err := utils.Wrap(d.conn.Where("to_user_id = ?", d.loginUserID).Order("create_time DESC").Find(&friendRequestList).Error, "GetRecvFriendApplication failed")

	var transfer []*model_struct.LocalFriendRequest
	for _, v := range friendRequestList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetRecvFriendApplication failed")
}

func (d *DataBase) GetSendFriendApplication() ([]*model_struct.LocalFriendRequest, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()
	var friendRequestList []model_struct.LocalFriendRequest
	err := utils.Wrap(d.conn.Where("from_user_id = ?", d.loginUserID).Order("create_time DESC").Find(&friendRequestList).Error, "GetSendFriendApplication failed")

	var transfer []*model_struct.LocalFriendRequest
	for _, v := range friendRequestList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetSendFriendApplication failed")
}

func (d *DataBase) GetFriendApplicationByBothID(fromUserID, toUserID string) (*model_struct.LocalFriendRequest, error) {
	d.friendMtx.Lock()
	defer d.friendMtx.Unlock()

	var friendRequest model_struct.LocalFriendRequest
	err := utils.Wrap(d.conn.Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).Take(&friendRequest).Error, "GetFriendApplicationByBothID failed")

	return &friendRequest, utils.Wrap(err, "GetFriendApplicationByBothID failed")
}
