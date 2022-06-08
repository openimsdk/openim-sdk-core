package db

import (
	"fmt"
	"open_im_sdk/pkg/utils"
	"strings"
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
	err := d.getDepartmentList(&departmentList, departmentID)
	return departmentList, err
}

func (d *DataBase) getDepartmentList(departmentList *[]*LocalDepartment, departmentID string) error {
	if len(*departmentList) == 0 {
		department, err := d.GetDepartmentInfo(departmentID)
		if err != nil {
			return err
		}
		*departmentList = append(*departmentList, department)
	}
	department, err := d.getParentDepartment(departmentID)
	if err != nil {
		return utils.Wrap(err, "getParentDepartment failed")
	}
	if department.DepartmentID != "" {
		*departmentList = append([]*LocalDepartment{&department}, *departmentList...)
		err := d.getDepartmentList(departmentList, department.DepartmentID)
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

func (d *DataBase) SearchDepartmentMember(keyWord string, isSearchUserName, isSearchEmail, isSearchMobile, isSearchPosition, isSearchTelephone, isSearchUserEnglishName, isSearchUserID bool, offset, count int) ([]*SearchDepartmentMemberResult, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var departmentMemberList []*SearchDepartmentMemberResult
	likeCondition := fmt.Sprintf("%%%s%%", keyWord)
	var likeConditions []interface{}
	var whereConditions []string
	if isSearchEmail {
		whereConditions = append(whereConditions, "email LIKE ?")
		likeConditions = append(likeConditions, likeCondition)
	}
	if isSearchMobile {
		whereConditions = append(whereConditions, "mobile LIKE ?")
		likeConditions = append(likeConditions, likeCondition)
	}
	if isSearchPosition {
		whereConditions = append(whereConditions, "position LIKE ?")
		likeConditions = append(likeConditions, likeCondition)
	}
	if isSearchTelephone {
		whereConditions = append(whereConditions, "telephone LIKE ?")
		likeConditions = append(likeConditions, likeCondition)

	}
	if isSearchUserEnglishName {
		whereConditions = append(whereConditions, "english_name LIKE ?")
		likeConditions = append(likeConditions, likeCondition)
	}
	if isSearchUserID {
		whereConditions = append(whereConditions, "user_id LIKE ?")
		likeConditions = append(likeConditions, likeCondition)
	}
	if isSearchUserName {
		whereConditions = append(whereConditions, "nickname LIKE ?")
		likeConditions = append(likeConditions, likeCondition)
	}
	if len(whereConditions) == 0 {
		whereConditions = append(whereConditions, "nickname LIKE ?")
		likeConditions = append(likeConditions, likeCondition)
	}
	err := d.conn.Model(&LocalDepartmentMember{}).
		Select([]string{"local_departments.name, local_department_members.*"}).
		//Where("nickname LIKE ? or english_name LIKE ? or mobile LIKE ? or telephone LIKE ? or email LIKE ? or position LIKE ?",
		Where(strings.Join(whereConditions, " or "),
			likeConditions...).
		Joins("left join local_departments on local_departments.department_id = local_department_members.department_id").
		Offset(offset).Limit(count).
		Scan(&departmentMemberList).Error
	return departmentMemberList, utils.Wrap(err, "")
}

func (d *DataBase) SearchDepartment(keyWord string, offset, count int) ([]*LocalDepartment, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var departmentMemberList []*LocalDepartment
	likeCondition := fmt.Sprintf("%%%s%%", keyWord)
	err := d.conn.Model(&LocalDepartment{}).Where("name LIKE ? and department_id != 0", likeCondition).Offset(offset).Limit(count).Find(&departmentMemberList).Error
	return departmentMemberList, utils.Wrap(err, "")
}
