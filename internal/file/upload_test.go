package file

import (
	"context"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/sdk_struct"
	"path/filepath"
	"testing"
)

func TestName(t *testing.T) {
	userID := `11111112`
	ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: userID,
		Token:  `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiIxMTExMTExMiIsIlBsYXRmb3JtSUQiOjEsImV4cCI6MTY5Njc2MjY5NywibmJmIjoxNjg4OTg2Mzk3LCJpYXQiOjE2ODg5ODY2OTd9.Cl7QPIRHWeUmF1FmC0Z8Kk3AFO8WbHrW6N2GkuG2hFc`,
		IMConfig: sdk_struct.IMConfig{
			ApiAddr: "http://125.124.195.201:10002",
		},
	})
	ctx = ccontext.WithOperationID(ctx, `test`)
	f := NewFile(nil, userID)

	//resp, err := f.accessURL(ctx, &third.AccessURLReq{
	//	Name: `11111112/test.png`,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(resp)

	path := `C:\Users\Admin\Desktop\test.png`
	resp, err := f.UploadFile(ctx, &UploadFileReq{
		Filepath: path,
		Name:     filepath.Base(path),
		Cause:    "test",
	}, nil)
	if err != nil {
		t.Logf("%+v\n", err)
		return
	}
	t.Logf("%+v\n", resp)
}
