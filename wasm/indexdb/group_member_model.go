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
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/exec"
)

type LocalGroupMember struct {
}

func NewLocalGroupMember() *LocalGroupMember {
	return &LocalGroupMember{}
}

func (i *LocalGroupMember) GetGroupMemberInfoByGroupIDUserID(ctx context.Context, groupID, userID string) (*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID, userID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return &temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetAllGroupMemberList(ctx context.Context) ([]model_struct.LocalGroupMember, error) {
	member, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetAllGroupMemberUserIDList(ctx context.Context) ([]model_struct.LocalGroupMember, error) {
	member, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberCount(ctx context.Context, groupID string) (int32, error) {
	count, err := exec.Exec(groupID)
	if err != nil {
		return 0, err
	}
	if v, ok := count.(float64); ok {
		return int32(v), nil
	}
	return 0, exec.ErrType
}

func (i *LocalGroupMember) GetGroupSomeMemberInfo(ctx context.Context, groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID, utils.StructToJsonString(userIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupAdminID(ctx context.Context, groupID string) ([]string, error) {
	IDList, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	}
	if v, ok := IDList.(string); ok {
		var temp []string
		err := utils.JsonStringToStruct(v, &temp)
		if err != nil {
			return nil, err
		}
		return temp, nil
	}
	return nil, exec.ErrType
}

func (i *LocalGroupMember) GetGroupMemberListByGroupID(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberListSplit(ctx context.Context, groupID string, filter int32, offset, count int) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID, filter, offset, count)
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberOwnerAndAdminDB(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberOwner(ctx context.Context, groupID string) (*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return &temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberListSplitByJoinTimeFilter(ctx context.Context, groupID string, offset, count int, joinTimeBegin, joinTimeEnd int64, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID, offset, count, joinTimeBegin, joinTimeEnd, utils.StructToJsonString(userIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupOwnerAndAdminByGroupID(ctx context.Context, groupID string) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberUIDListByGroupID(ctx context.Context, groupID string) (result []string, err error) {
	IDList, err := exec.Exec(groupID)
	if err != nil {
		return nil, err
	}
	if v, ok := IDList.(string); ok {
		err := utils.JsonStringToStruct(v, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, exec.ErrType
}

func (i *LocalGroupMember) InsertGroupMember(ctx context.Context, groupMember *model_struct.LocalGroupMember) error {
	_, err := exec.Exec(utils.StructToJsonString(groupMember))
	return err
}

func (i *LocalGroupMember) BatchInsertGroupMember(ctx context.Context, groupMemberList []*model_struct.LocalGroupMember) error {
	_, err := exec.Exec(utils.StructToJsonString(groupMemberList))
	return err
}

func (i *LocalGroupMember) DeleteGroupMember(ctx context.Context, groupID, userID string) error {
	_, err := exec.Exec(groupID, userID)
	return err
}

func (i *LocalGroupMember) DeleteGroupAllMembers(ctx context.Context, groupID string) error {
	_, err := exec.Exec(groupID)
	return err
}

func (i *LocalGroupMember) UpdateGroupMember(ctx context.Context, groupMember *model_struct.LocalGroupMember) error {
	_, err := exec.Exec(utils.StructToJsonString(groupMember))
	return err
}

func (i *LocalGroupMember) UpdateGroupMemberField(ctx context.Context, groupID, userID string, args map[string]interface{}) error {
	_, err := exec.Exec(groupID, userID, utils.StructToJsonString(args))
	return err
}

func (i *LocalGroupMember) GetGroupMemberInfoIfOwnerOrAdmin(ctx context.Context) ([]*model_struct.LocalGroupMember, error) {
	member, err := exec.Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) SearchGroupMembersDB(ctx context.Context, keyword string, groupID string, isSearchMemberNickname, isSearchUserID bool, offset, count int) (result []*model_struct.LocalGroupMember, err error) {
	member, err := exec.Exec(keyword, groupID, isSearchMemberNickname, isSearchUserID, offset, count)
	if err != nil {
		return nil, err
	} else {
		if v, ok := member.(string); ok {
			var temp []*model_struct.LocalGroupMember
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			return temp, err
		} else {
			return nil, exec.ErrType
		}
	}
}

func (i *LocalGroupMember) GetUserJoinedGroupIDs(ctx context.Context, userID string) (result []string, err error) {
	IDList, err := exec.Exec(userID)
	if err != nil {
		return nil, err
	}
	if v, ok := IDList.(string); ok {
		err := utils.JsonStringToStruct(v, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, exec.ErrType
}
