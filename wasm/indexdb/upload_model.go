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

//go:build js && wasm
// +build js,wasm

package indexdb

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
)

type LocalUpload struct{}

func NewLocalUpload() *LocalUpload {
	return &LocalUpload{}
}

func (i *LocalUpload) GetUpload(ctx context.Context, partHash string) (*model_struct.Upload, error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalUpload) InsertUpload(ctx context.Context, upload *model_struct.Upload) error {
	//TODO implement me
	panic("implement me")
}

func (i *LocalUpload) DeleteUpload(ctx context.Context, partHash string) error {
	//TODO implement me
	panic("implement me")
}

func (i *LocalUpload) DeleteExpireUpload(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (i *LocalUpload) GetUploadPart(ctx context.Context, partHash string) ([]int32, error) {
	//TODO implement me
	panic("implement me")
}

func (i *LocalUpload) SetUploadPartPush(ctx context.Context, partHash string, index []int32) error {
	//TODO implement me
	panic("implement me")
}
