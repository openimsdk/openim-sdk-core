package testv3

import (
	"open_im_sdk/pkg/utils"
	"open_im_sdk/testv3/funcation"
	"strconv"
	"testing"
	"time"
)

func Test_RegisterOne(t *testing.T) {
	uid := "bantanger"
	nickname := "bantanger"
	faceUrl := ""
	register, err := funcation.RegisterOne(uid, nickname, faceUrl)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(register)
}

func Test_RegisterBatch(t *testing.T) {
	count := 100
	var users []*funcation.Users
	for i := 0; i < count; i++ {
		users[i].Uid = funcation.GenUid(i, "register_test_"+utils.Int64ToString(time.Now().Unix()))
		users[i].Nickname = "register_test_" + strconv.FormatInt(int64(i), 10)
		users[i].FaceUrl = ""
	}
	success, fail := funcation.RegisterBatch(users)
	t.Log(success)
	t.Log(fail)
}
