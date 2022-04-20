package organization

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
)

func (o *Organization) SetListener(callback open_im_sdk_callback.OnOrganizationListener) {
	if callback == nil {
		return
	}
	o.listener = callback
}

func (o *Organization) GetSubDepartment(callback open_im_sdk_callback.Base, departmentID string, offset, count int, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", departmentID, offset, count)
		result := o.getSubDepartment(callback, departmentID, offset, count, operationID)
		resp := utils.StructToJsonStringDefault(result)
		callback.OnSuccess(resp)
		log.NewInfo(operationID, fName, " callback: ", resp)
	}()
}

func (o *Organization) GetDepartmentMember(callback open_im_sdk_callback.Base, departmentID string, offset, count int, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", departmentID, offset, count)
		result := o.getDepartmentMember(callback, departmentID, offset, count, operationID)
		resp := utils.StructToJsonStringDefault(result)
		callback.OnSuccess(resp)
		log.NewInfo(operationID, fName, " callback: ", resp)
	}()
}

func (o *Organization) GetUserInDepartment(callback open_im_sdk_callback.Base, userID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", userID)
		result := o.getUserInDepartment(callback, userID, operationID)
		resp := utils.StructToJsonStringDefault(result)
		callback.OnSuccess(resp)
		log.NewInfo(operationID, fName, " callback: ", resp)
	}()
}

func (o *Organization) GetDepartmentMemberAndSubDepartment(callback open_im_sdk_callback.Base, departmentID string, departmentOffset, departmentCount, memberOffset, memberCount int, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", departmentID, departmentOffset, departmentCount, memberOffset, memberCount)
		result := o.getDepartmentMemberAndSubDepartment(callback, departmentID, departmentOffset, departmentCount, memberOffset, memberCount, operationID)
		resp := utils.StructToJsonStringDefault(result)
		callback.OnSuccess(resp)
		log.NewInfo(operationID, fName, " callback: ", resp)
	}()
}
