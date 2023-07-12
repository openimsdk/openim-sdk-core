// @Author BanTanger 2023/7/10 12:30:00
package testv3

import (
	"open_im_sdk/testv3/funcation"
	"testing"
	"time"
)

func Test_LoginOne(t *testing.T) {
	uid := "6506148011"
	res := funcation.LoginOne(uid)
	time.Sleep(1000 * time.Second)
	if res != true {
		t.Errorf("uid [%v] login expected be successful, but fail got", uid)
	}
	t.Logf("uid [%v] login successfully", uid)
}

func Test_LoginBatch(t *testing.T) {
	count := 100
	userIDList := funcation.AllUserID
	funcation.LoginBatch(userIDList[:count])
}
