// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package file

import (
	"context"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/db"
	"open_im_sdk/sdk_struct"
	"path/filepath"
	"testing"
	"time"
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

	database, err := db.NewDataBase(ctx, userID, `C:\Users\Admin\Desktop\test`, 6)
	if err != nil {
		panic(err)
	}
	f := NewFile(database, userID)

	go func() {
		path := `C:\Users\Admin\Desktop\test`
		path = filepath.Join(path, `1.png`)
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
	}()

	go func() {
		path := `C:\Users\Admin\Desktop\test`
		path = filepath.Join(path, `2.png`)
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
	}()

	go func() {
		path := `C:\Users\Admin\Desktop\test`
		path = filepath.Join(path, `3.png`)
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
	}()

	time.Sleep(time.Second * 100)

}
