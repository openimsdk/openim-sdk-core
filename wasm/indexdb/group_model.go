//go:build js && wasm
// +build js,wasm

package indexdb

import "context"

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

type LocalGroups struct{}

func NewLocalGroups() *LocalGroups {
	return &LocalGroups{}
}

func (i *LocalGroups) InsertGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	_, err := Exec(utils.StructToJsonString(groupInfo))
	return err
}

func (i *LocalGroups) DeleteGroup(ctx context.Context, groupID string) error {
	_, err := Exec(groupID)
	return err
}

// 该函数需要全更新
func (i *LocalGroups) UpdateGroup(ctx context.Context, groupInfo *model_struct.LocalGroup) error {
	_, err := Exec(groupInfo.GroupID, utils.StructToJsonString(groupInfo))
	return err
}

func (i *LocalGroups) GetJoinedGroupListDB(ctx context.Context) (result []*model_struct.LocalGroup, err error) {
	gList, err := Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalGroup
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalGroups) GetGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	c, err := Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := c.(string); ok {
			result := model_struct.LocalGroup{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalGroups) GetAllGroupInfoByGroupIDOrGroupName(ctx context.Context, keyword string, isSearchGroupID bool, isSearchGroupName bool) (result []*model_struct.LocalGroup, err error) {
	gList, err := Exec(keyword, isSearchGroupID, isSearchGroupName)
	if err != nil {
		return nil, err
	} else {
		if v, ok := gList.(string); ok {
			var temp []model_struct.LocalGroup
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalGroups) AddMemberCount(ctx context.Context, groupID string) error {
	_, err := Exec(groupID)
	return err
}

func (i *LocalGroups) SubtractMemberCount(ctx context.Context, groupID string) error {
	_, err := Exec(groupID)
	return err
}
