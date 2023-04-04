package open_im_sdk

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
)

func GetSubDepartment(callback open_im_sdk_callback.Base, operationID string, departmentID string, offset, count int) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Organization().GetSubDepartment(callback, departmentID, offset, count, operationID)
}

func GetDepartmentMember(callback open_im_sdk_callback.Base, operationID string, departmentID string, offset, count int) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Organization().GetDepartmentMember(callback, departmentID, offset, count, operationID)
}

func GetUserInDepartment(callback open_im_sdk_callback.Base, operationID string, userID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Organization().GetUserInDepartment(callback, userID, operationID)
}

func GetDepartmentMemberAndSubDepartment(callback open_im_sdk_callback.Base, operationID string, departmentID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Organization().GetDepartmentMemberAndSubDepartment(callback, departmentID, operationID)
}

func GetParentDepartmentList(callback open_im_sdk_callback.Base, operationID string, departmentID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Organization().GetParentDepartmentList(callback, departmentID, operationID)
}

func GetDepartmentInfo(callback open_im_sdk_callback.Base, operationID string, departmentID string) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Organization().GetDepartmentInfo(callback, departmentID, operationID)
}

func SearchOrganization(callback open_im_sdk_callback.Base, operationID string, input string, offset, count int) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Organization().SearchOrganization(callback, input, offset, count, operationID)
}
