package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

type LocalGroupRequest struct {
}

func (i *LocalGroupRequest) InsertGroupRequest(groupRequest *model_struct.LocalGroupRequest) error {
	_, err := Exec(utils.StructToJsonString(groupRequest))
	return err
}

func (i *LocalGroupRequest) DeleteGroupRequest(groupID, userID string) error {
	_, err := Exec(groupID, userID)
	return err
}

func (i *LocalGroupRequest) UpdateGroupRequest(groupRequest *model_struct.LocalGroupRequest) error {
	_, err := Exec(utils.StructToJsonString(groupRequest))
	return err
}

func (i *LocalGroupRequest) GetSendGroupApplication() ([]*model_struct.LocalGroupRequest, error) {
	result, err := Exec()
	if err != nil {
		return nil, err
	}
	if v, ok := result.(string); ok {
		var request []*model_struct.LocalGroupRequest
		if err := utils.JsonStringToStruct(v, &request); err != nil {
			return nil, err
		}
		return request, nil
	} else {
		return nil, ErrType
	}
}

func (i IndexDB) InsertAdminGroupRequest(groupRequest *model_struct.LocalAdminGroupRequest) error {
	_, err := Exec(utils.StructToJsonString(groupRequest))
	return err
}

func (i IndexDB) DeleteAdminGroupRequest(groupID, userID string) error {
	_, err := Exec(groupID, userID)
	return err
}

func (i IndexDB) UpdateAdminGroupRequest(groupRequest *model_struct.LocalAdminGroupRequest) error {
	_, err := Exec(utils.StructToJsonString(groupRequest))
	return err
}

func (i IndexDB) GetAdminGroupApplication() ([]*model_struct.LocalAdminGroupRequest, error) {
	result, err := Exec()
	if err != nil {
		return nil, err
	}
	if v, ok := result.(string); ok {
		var request []*model_struct.LocalAdminGroupRequest
		if err := utils.JsonStringToStruct(v, &request); err != nil {
			return nil, err
		}
		return request, nil
	} else {
		return nil, ErrType
	}
}
