package db

import (
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertFriendRequest(friendRequest *LocalFriendRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(friendRequest).Error, "InsertFriendRequest failed")
	//u := d.conn.Model(friendRequest).Updates(args)
	//if u.RowsAffected != 0 {
	//	return nil
	//}

}

func (d *DataBase) DeleteFriendRequestBothUserID(fromUserID, toUserID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Where("from_user_id=? and to_user_id=?", fromUserID, toUserID).Delete(&LocalFriendRequest{}).Error, "DeleteFriendRequestBothUserID failed")
}

func (d *DataBase) UpdateFriendRequest(friendRequest *LocalFriendRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Model(friendRequest).Select("*").Updates(*friendRequest).Error, "")

	//
	//t := d.conn.Updates(friendRequest)
	//if t.RowsAffected == 0 {
	//	return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	//}
	//return utils.Wrap(t.Error, "UpdateFriendRequest failed")
}

func (d *DataBase) GetRecvFriendApplication() ([]*LocalFriendRequest, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var friendRequestList []LocalFriendRequest
	err := utils.Wrap(d.conn.Where("to_user_id = ?", d.loginUserID).Find(&friendRequestList).Error, "GetRecvFriendApplication failed")

	var transfer []*LocalFriendRequest
	for _, v := range friendRequestList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetRecvFriendApplication failed")
}

func (d *DataBase) GetSendFriendApplication() ([]*LocalFriendRequest, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var friendRequestList []LocalFriendRequest
	err := utils.Wrap(d.conn.Where("from_user_id = ?", d.loginUserID).Find(&friendRequestList).Error, "GetSendFriendApplication failed")

	var transfer []*LocalFriendRequest
	for _, v := range friendRequestList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetSendFriendApplication failed")
}

func (d *DataBase) GetFriendApplicationByBothID(fromUserID, toUserID string) (*LocalFriendRequest, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	var friendRequest LocalFriendRequest
	err := utils.Wrap(d.conn.Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).Take(&friendRequest).Error, "GetFriendApplicationByBothID failed")

	return &friendRequest, utils.Wrap(err, "GetFriendApplicationByBothID failed")
}
