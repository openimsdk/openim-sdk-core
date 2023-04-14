//go:build !js
// +build !js

package db

import "context"

import (
	"errors"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetDepartmentMemberListByDepartmentID(ctx context.Context, departmentID string, args ...int) ([]*model_struct.LocalDepartmentMember, error) {
	d.departmentMtx.RLock()
	defer d.departmentMtx.RUnlock()
	var departmentMemberList []model_struct.LocalDepartmentMember
	var err error
	sql := d.conn.Where("department_id = ? ", departmentID).Order("order_member DESC")
	if len(args) == 2 {
		offset := args[0]
		count := args[1]
		err = sql.Offset(offset).Limit(count).Find(&departmentMemberList).Error
	} else {
		err = sql.Find(&departmentMemberList).Error
	}
	var transfer []*model_struct.LocalDepartmentMember
	for _, v := range departmentMemberList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) GetAllDepartmentMemberList(ctx context.Context) ([]*model_struct.LocalDepartmentMember, error) {
	d.departmentMtx.RLock()
	defer d.departmentMtx.RUnlock()
	var departmentMemberList []model_struct.LocalDepartmentMember
	err := d.conn.Find(&departmentMemberList).Error

	var transfer []*model_struct.LocalDepartmentMember
	for _, v := range departmentMemberList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) InsertDepartmentMember(ctx context.Context, departmentMember *model_struct.LocalDepartmentMember) error {
	d.departmentMtx.RLock()
	defer d.departmentMtx.RUnlock()
	return utils.Wrap(d.conn.Create(departmentMember).Error, "InsertDepartmentMember failed")
}

func (d *DataBase) BatchInsertDepartmentMember(ctx context.Context, departmentMemberList []*model_struct.LocalDepartmentMember) error {
	d.departmentMtx.RLock()
	defer d.departmentMtx.RUnlock()
	if departmentMemberList == nil {
		return errors.New("nil")
	}
	return utils.Wrap(d.conn.Create(departmentMemberList).Error, "BatchInsertDepartmentMember failed")
}

//func (d *DataBase) BatchInsertDepartmentMember(ctx context.Context, departmentMember *model_struct.LocalDepartmentMember) error {
//	d.mRWMutex.Lock()
//	defer d.mRWMutex.Unlock()
//	return utils.Wrap(d.conn.Create(departmentMember).Error, "InsertDepartmentMember failed")
//}

func (d *DataBase) UpdateDepartmentMember(ctx context.Context, departmentMember *model_struct.LocalDepartmentMember) error {
	d.departmentMtx.RLock()
	defer d.departmentMtx.RUnlock()
	return utils.Wrap(d.conn.Model(departmentMember).Select("*").Updates(*departmentMember).Error, "UpdateDepartmentMember failed")
}

func (d *DataBase) DeleteDepartmentMember(ctx context.Context, departmentID string, userID string) error {
	d.departmentMtx.RLock()
	defer d.departmentMtx.RUnlock()
	//local := LocalDepartmentMember{DepartmentID: departmentID, UserID: userID}
	return utils.Wrap(d.conn.Where("department_id = ? and user_id = ?", departmentID, userID).Delete(&model_struct.LocalDepartmentMember{}).Error, "DeleteDepartmentMember failed")
}

func (d *DataBase) GetDepartmentMemberListByUserID(ctx context.Context, userID string) ([]*model_struct.LocalDepartmentMember, error) {
	d.departmentMtx.RLock()
	defer d.departmentMtx.RUnlock()
	var departmentMemberList []model_struct.LocalDepartmentMember
	err := d.conn.Where("user_id = ? ", userID).Find(&departmentMemberList).Error
	var transfer []*model_struct.LocalDepartmentMember
	for _, v := range departmentMemberList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, utils.GetSelfFuncName()+" failed")
}
