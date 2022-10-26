package test

import (
	"fmt"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
)

type userBaseCallback struct {
	OperationID string
	CallName    string
}

func (t userBaseCallback) OnSuccess(data string) {
	fmt.Println("=======================\n", t.OperationID, t.CallName, utils.GetSelfFuncName(), data)
}

func (t userBaseCallback) OnError(errCode int32, errMsg string) {
	fmt.Println(t.OperationID, t.CallName, utils.GetSelfFuncName(), errCode, errMsg)
}

func DoTestSqlTest() {
	var test userBaseCallback
	test.OperationID = utils.OperationIDGenerator()
	test.CallName = utils.GetSelfFuncName()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ")
	open_im_sdk.SqlTest(test)
}
