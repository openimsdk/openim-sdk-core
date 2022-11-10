package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

type LocalGroupMember struct {
}

func (i *LocalGroupMember) GetGroupMemberInfoByGroupIDUserID(groupID, userID string) (*model_struct.LocalGroupMember, error) {
	member, err := Exec(groupID, userID)
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
			return nil, ErrType
		}
	}
}

func (i *LocalGroupMember) GetAllGroupMemberList() ([]model_struct.LocalGroupMember, error) {
	member, err := Exec()
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
			return nil, ErrType
		}
	}
}

func (i *LocalGroupMember) GetAllGroupMemberUserIDList() ([]model_struct.LocalGroupMember, error) {
	member, err := Exec()
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
			return nil, ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberCount(groupID string) (uint32, error) {
	count, err := Exec(groupID)
	if err != nil {
		return 0, err
	}
	if v, ok := count.(float64); ok {
		return uint32(v), nil
	}
	return 0, ErrType
}

func (i *LocalGroupMember) GetGroupSomeMemberInfo(groupID string, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	member, err := Exec(groupID, utils.StructToJsonString(userIDList))
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
			return nil, ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupAdminID(groupID string) ([]string, error) {
	IDList, err := Exec(groupID)
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

func (i *LocalGroupMember) GetGroupMemberListByGroupID(groupID string) ([]*model_struct.LocalGroupMember, error) {
	member, err := Exec(groupID)
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
			return nil, ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberListSplit(groupID string, filter int32, offset, count int) ([]*model_struct.LocalGroupMember, error) {
	member, err := Exec(groupID, filter, offset, count)
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
			return nil, ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberOwnerAndAdmin(groupID string) ([]*model_struct.LocalGroupMember, error) {
	member, err := Exec(groupID)
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
			return nil, ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberOwner(groupID string) (*model_struct.LocalGroupMember, error) {
	member, err := Exec(groupID)
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
			return nil, ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberListSplitByJoinTimeFilter(groupID string, offset, count int, joinTimeBegin, joinTimeEnd int64, userIDList []string) ([]*model_struct.LocalGroupMember, error) {
	member, err := Exec(groupID, offset, count, joinTimeBegin, joinTimeEnd, utils.StructToJsonString(userIDList))
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
			return nil, ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupOwnerAndAdminByGroupID(groupID string) ([]*model_struct.LocalGroupMember, error) {
	member, err := Exec(groupID)
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
			return nil, ErrType
		}
	}
}

func (i *LocalGroupMember) GetGroupMemberUIDListByGroupID(groupID string) (result []string, err error) {
	IDList, err := Exec(groupID)
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
	return nil, ErrType
}

func (i *LocalGroupMember) InsertGroupMember(groupMember *model_struct.LocalGroupMember) error {
	_, err := Exec(utils.StructToJsonString(groupMember))
	return err
}

func (i *LocalGroupMember) BatchInsertGroupMember(groupMemberList []*model_struct.LocalGroupMember) error {
	_, err := Exec(utils.StructToJsonString(groupMemberList))
	return err
}

func (i *LocalGroupMember) DeleteGroupMember(groupID, userID string) error {
	_, err := Exec(groupID, userID)
	return err
}

func (i *LocalGroupMember) DeleteGroupAllMembers(groupID string) error {
	_, err := Exec(groupID)
	return err
}

func (i *LocalGroupMember) UpdateGroupMember(groupMember *model_struct.LocalGroupMember) error {
	_, err := Exec(utils.StructToJsonString(groupMember))
	return err
}

func (i *LocalGroupMember) UpdateGroupMemberField(groupID, userID string, args map[string]interface{}) error {
	_, err := Exec(groupID, userID, utils.StructToJsonString(args))
	return err
}

func (i *LocalGroupMember) GetGroupMemberInfoIfOwnerOrAdmin() ([]*model_struct.LocalGroupMember, error) {
	member, err := Exec()
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
			return nil, ErrType
		}
	}
}

func (i *LocalGroupMember) SearchGroupMembersDB(keyword string, groupID string, isSearchMemberNickname, isSearchUserID bool, offset, count int) (result []*model_struct.LocalGroupMember, err error) {
	member, err := Exec(keyword, groupID, isSearchMemberNickname, isSearchUserID, offset, count)
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
			return nil, ErrType
		}
	}
}
