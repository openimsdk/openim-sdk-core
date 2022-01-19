package db

import (
	"errors"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertGroupRequest(groupRequest *LocalGroupRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(groupRequest).Error, "InsertGroupRequest failed")

}
func (d *DataBase) DeleteGroupRequest(groupID, userID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Where("group_id=? and user_id=?", groupID, userID).Delete(&LocalGroupRequest{}).Error, "DeleteGroupRequest failed")
}
func (d *DataBase) UpdateGroupRequest(groupRequest *LocalGroupRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Updates(groupRequest)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "_updateGroupRequest failed")
}

//
//func (d *DataBase) GetAdminGroupApplication() ([]*LocalGroupRequest, error) {
//	ownerAdminList, err := d.GetGroupMenberInfoIfOwnerOrAdmin()
//	if err != nil {
//		return nil, utils.Wrap(err, "")
//	}
//
//	//fmt.Println("ownerAdminList ", ownerAdminList)
//	d.mRWMutex.Lock()
//	defer d.mRWMutex.Unlock()
//	var transfer []*LocalGroupRequest
//	for _, v := range ownerAdminList {
//		var groupRequest LocalGroupRequest
//		f := d.conn.Where("group_id = ?", v.GroupID).Find(&groupRequest)
//		err := f.Error
//
//		if err != nil {
//			continue
//		}
//		if f.RowsAffected != 0 {
//			transfer = append(transfer, &groupRequest)
//		}
//	}
//	return transfer, utils.Wrap(err, "GetRecvGroupApplication failed ")
//}

func (d *DataBase) GetSelfGroupApplication() ([]LocalGroupRequest, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var groupRequestList []LocalGroupRequest
	return groupRequestList, utils.Wrap(d.conn.Where("user_id = ?", d.loginUserID).Find(&groupRequestList).Error, "GetSendGroupApplication failed")
}
