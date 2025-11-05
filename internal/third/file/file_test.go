package file

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"path/filepath"
	"testing"
)

func TestUpload(t *testing.T) {
	conf := &ccontext.GlobalConfig{
		UserID: `4931176757`,
		Token:  `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI0OTMxMTc2NzU3IiwiUGxhdGZvcm1JRCI6MSwiZXhwIjoxNzA3MTE0MjIyLCJuYmYiOjE2OTkzMzc5MjIsImlhdCI6MTY5OTMzODIyMn0.AyNvrMGEdXD5rkvn7ZLHCNs-lNbDCb2otn97yLXia5Y`,
		IMConfig: sdk_struct.IMConfig{
			ApiAddr: `http://203.56.175.233:10002`,
		},
	}
	ctx := ccontext.WithInfo(context.WithValue(context.Background(), "operationID", "OP123456"), conf)
	f := NewFile(nil, conf.UserID)

	fp := `C:\Users\openIM\Desktop\protoc.zip`

	resp, err := f.UploadFile(ctx, &UploadFileReq{
		Filepath: fp,
		Name:     filepath.Base(fp),
		Cause:    "test",
	}, nil)
	if err != nil {
		t.Fatal("failed", err)
	}
	t.Log("success", resp.URL)

}
