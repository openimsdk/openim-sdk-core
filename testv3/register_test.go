package testv3

import (
	"open_im_sdk/testv3/funcation"
	"testing"
)

func Test_Register(t *testing.T) {
	uid := "123456"
	nickname := "bantanger"
	faceurl := ""
	register, err := funcation.Register(uid, nickname, faceurl)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(register)
}
