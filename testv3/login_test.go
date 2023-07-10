package testv3

import (
	"open_im_sdk/testv3/funcation"
	"testing"
)

func Test_LoginOne(t *testing.T) {
	idx := 0
	uid := "bantanger"
	funcation.LoginOne(idx, uid)
	t.Log("uid [ " + uid + " ] login successfully")
}

func Test_LoginBatch(t *testing.T) {
	count := 100
	var uidList []string
	for i := 0; i < count; i++ {
		user := funcation.AllLoginMgr[i]
		uidList = append(uidList, user.UserID)
	}
	funcation.LoginBatch(uidList)
}
