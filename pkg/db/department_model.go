package db

import (
	"fmt"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetSubDepartmentList(departmentID string, args ...int) ([]*LocalDepartment, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var departmentList []LocalDepartment
	var err error
	sql := d.conn.Where("parent_id = ? ", departmentID).Order("order_department DESC")
	if len(args) == 2 {
		offset := args[0]
		count := args[1]
		err = sql.Offset(offset).Limit(count).Find(&departmentList).Error
	} else {
		err = sql.Find(&departmentList).Error
	}
	var transfer []*LocalDepartment
	for _, v := range departmentList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) InsertDepartment(department *LocalDepartment) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(department).Error, "InsertDepartment failed")
}

func (d *DataBase) UpdateDepartment(department *LocalDepartment) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Model(department).Select("*").Updates(*department).Error, "UpdateDepartment failed")
}

func (d *DataBase) DeleteDepartment(departmentID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	local := LocalDepartment{DepartmentID: departmentID}
	return utils.Wrap(d.conn.Delete(&local).Error, "DeleteDepartment failed")
}

func (d *DataBase) GetDepartmentInfo(departmentID string) (*LocalDepartment, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var local LocalDepartment
	return &local, utils.Wrap(d.conn.Where("department_id=?", departmentID).First(&local).Error, "GetDepartmentInfo failed")
}

func (d *DataBase) GetAllDepartmentList() ([]*LocalDepartment, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var departmentList []LocalDepartment
	d.conn.Debug()
	//	err := d.conn.Order("order DESC").Find(&departmentList).Error
	err := d.conn.Order("order_department DESC").Find(&departmentList).Error
	var transfer []*LocalDepartment
	for _, v := range departmentList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, utils.GetSelfFuncName()+" failed")
}

func (d *DataBase) GetParentDepartmentList(departmentID string) ([]*LocalDepartment, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var departmentList []*LocalDepartment
	err := d.getParentDepartmentList(&departmentList, departmentID)
	return departmentList, err
}

func (d *DataBase) getParentDepartmentList(departmentList *[]*LocalDepartment, departmentID string) error {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	department, err := d.getParentDepartment(departmentID)
	if err != nil {
		return utils.Wrap(err, "getParentDepartment failed")
	}
	if department.DepartmentID != "" {
		*departmentList = append([]*LocalDepartment{&department}, *departmentList...)
		err := d.getParentDepartmentList(departmentList, department.DepartmentID)
		if err != nil {
			return utils.Wrap(err, "getParentDepartmentList failed")
		}
	}
	return nil
}

func (d *DataBase) getParentDepartment(departmentID string) (LocalDepartment, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var department LocalDepartment
	var parentID string
	d.conn.Model(&department).Where("department_id=?", departmentID).Pluck("parent_id", &parentID)
	err := d.conn.Where("department_id=?", parentID).Find(&department).Error
	return department, utils.Wrap(err, "getParentDepartment failed")
}

type SearchDepartmentMemberResult struct {
	LocalDepartmentMember
	DepartmentName string `gorm:"column:name;size:256" json:"departmentName"`
}

func (d *DataBase) SearchDepartmentMember(input string, offset, count int) ([]*SearchDepartmentMemberResult, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var departmentMemberList []*SearchDepartmentMemberResult
	condition := fmt.Sprintf("%%%s%%", input)
	err := d.conn.Model(&LocalDepartmentMember{}).
		Select([]string{"local_departments.name, local_department_members.*"}).
		Where("nickname LIKE ? or english_name LIKE ? or mobile LIKE ? or telephone LIKE ? or email LIKE ? or position LIKE ?",
			condition, condition, condition, condition, condition, condition).
		Joins("left join local_departments on local_departments.department_id = local_department_members.department_id").
		Offset(offset).Limit(count).
		Scan(&departmentMemberList).Error
	return departmentMemberList, utils.Wrap(err, "")
}

func (d *DataBase) SearchDepartment(input string, offset, count int) ([]*LocalDepartment, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var departmentMemberList []*LocalDepartment
	condition := fmt.Sprintf("%%%s%%", input)
	err := d.conn.Model(&LocalDepartment{}).Where("name LIKE ? and department_id != 0", condition).Offset(offset).Limit(count).Find(&departmentMemberList).Error
	return departmentMemberList, utils.Wrap(err, "")
}
