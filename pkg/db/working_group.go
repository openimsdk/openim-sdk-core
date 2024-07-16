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

//go:build !js
// +build !js

package db

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) GetJoinedWorkingGroupIDList(ctx context.Context) ([]string, error) {
	groupList, err := d.GetJoinedGroupListDB(ctx)
	if err != nil {
		return nil, errs.WrapMsg(err, "")
	}
	var groupIDList []string
	for _, v := range groupList {
		if v.GroupType == constant.WorkingGroup {
			groupIDList = append(groupIDList, v.GroupID)
		}
	}
	return groupIDList, nil
}

func (d *DataBase) GetJoinedWorkingGroupList(ctx context.Context) ([]*model_struct.LocalGroup, error) {
	groupList, err := d.GetJoinedGroupListDB(ctx)
	var transfer []*model_struct.LocalGroup
	for _, v := range groupList {
		if v.GroupType == constant.WorkingGroup {
			transfer = append(transfer, v)
		}
	}
	return transfer, errs.WrapMsg(err, "GetJoinedSuperGroupList failed ")
}
