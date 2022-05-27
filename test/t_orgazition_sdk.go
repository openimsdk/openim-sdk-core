package test

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
)

type orgBaseCallback struct {
	OperationID string
	CallName    string
}

func (t orgBaseCallback) OnSuccess(data string) {
	log.Info(t.OperationID, t.CallName, utils.GetSelfFuncName(), data)

}

func (t orgBaseCallback) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, t.CallName, utils.GetSelfFuncName(), errCode, errMsg)
}

type testOrganization struct {
	orgBaseCallback
}

func DoTestGetSubDepartment() {
	var test testOrganization
	test.OperationID = utils.OperationIDGenerator()
	test.CallName = utils.GetSelfFuncName()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ")
	open_im_sdk.GetSubDepartment(test, test.OperationID, "001", 0, 1)
}

func DoTestGetDepartmentMember() {
	var test testOrganization
	test.OperationID = utils.OperationIDGenerator()
	test.CallName = utils.GetSelfFuncName()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ")
	open_im_sdk.GetDepartmentMember(test, test.OperationID, "001", 0, 1)
}

func DoTestGetUserInDepartment() {
	var test testOrganization
	test.OperationID = utils.OperationIDGenerator()
	test.CallName = utils.GetSelfFuncName()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ")
	open_im_sdk.GetUserInDepartment(test, test.OperationID, "org_user_001")
}

func DoTestGetDepartmentMemberAndSubDepartment() {
	var test testOrganization
	test.OperationID = utils.OperationIDGenerator()
	test.CallName = utils.GetSelfFuncName()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ")
	open_im_sdk.GetDepartmentMemberAndSubDepartment(test, test.OperationID, "001", 0, 1, 0, 1)
}
