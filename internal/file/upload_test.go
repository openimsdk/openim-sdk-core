package file

import (
	"context"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/db"
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

	database, err := db.NewDataBase(ctx, userID, `C:\Users\Admin\Desktop\test`)
	if err != nil {
		panic(err)
	}
	f := NewFile(database, userID)

	path := `C:\Users\Admin\Desktop\test`
	path = filepath.Join(path, `go1.19.10.linux-amd64.tar.gz`)
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
