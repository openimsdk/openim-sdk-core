package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

type LocalSuperGroup struct{}

func (i *LocalSuperGroup) GetJoinedSuperGroupList() (result []*model_struct.LocalGroup, err error) {
	groupList, err := Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := groupList.(string); ok {
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
func (i *LocalSuperGroup) InsertSuperGroup(groupInfo *model_struct.LocalGroup) error {
	_, err := Exec(utils.StructToJsonString(groupInfo))
	return err
}
func (i *LocalSuperGroup) UpdateSuperGroup(g *model_struct.LocalGroup) error {
	_, err := Exec(g.GroupID, utils.StructToJsonString(g))
	return err
}

func (i *LocalSuperGroup) DeleteSuperGroup(groupID string) error {
	_, err := Exec(groupID)
	return err
}

func (i *LocalSuperGroup) DeleteAllSuperGroup() error {
	_, err := Exec()
	return err
}

func (i *LocalSuperGroup) GetSuperGroupInfoByGroupID(groupID string) (*model_struct.LocalGroup, error) {
	groupInfo, err := Exec(groupID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := groupInfo.(string); ok {
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

func (i *LocalSuperGroup) GetJoinedWorkingGroupIDList() ([]string, error) {
	IDList, err := Exec()
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
	return nil, ErrType
}

func (i *LocalSuperGroup) GetJoinedWorkingGroupList() (result []*model_struct.LocalGroup, err error) {
	groupList, err := Exec()
	if err != nil {
		return nil, err
	} else {
		if v, ok := groupList.(string); ok {
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
