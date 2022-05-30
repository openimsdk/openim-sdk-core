package db

import (
	"errors"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) InsertGroupRequest(groupRequest *model_struct.LocalGroupRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(groupRequest).Error, "InsertGroupRequest failed")

}
func (d *DataBase) DeleteGroupRequest(groupID, userID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Where("group_id=? and user_id=?", groupID, userID).Delete(&model_struct.LocalGroupRequest{}).Error, "DeleteGroupRequest failed")
}
func (d *DataBase) UpdateGroupRequest(groupRequest *model_struct.LocalGroupRequest) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()

	t := d.conn.Model(groupRequest).Select("*").Updates(*groupRequest)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "")
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

func (d *DataBase) GetSendGroupApplication() ([]*model_struct.LocalGroupRequest, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var groupRequestList []model_struct.LocalGroupRequest
	err := utils.Wrap(d.conn.Find(&groupRequestList).Error, "")
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var transfer []*model_struct.LocalGroupRequest
	for _, v := range groupRequestList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, nil
}
