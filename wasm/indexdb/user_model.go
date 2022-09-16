package indexdb

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

type LocalUsers struct {
}

func (l *LocalUsers) GetLoginUser(userID string) (*model_struct.LocalUser, error) {
	user, err := Exec(userID)
	if err != nil {
		return nil, err
	} else {
		result := model_struct.LocalUser{}
		err := utils.JsonStringToStruct(user, &result)
		if err != nil {
			return nil, err
		}
		return &result, err
	}
}

func (l *LocalUsers) UpdateLoginUser(user *model_struct.LocalUser) error {
	panic("implement me")

}
func (l *LocalUsers) UpdateLoginUserByMap(user *model_struct.LocalUser, args map[string]interface{}) error {
	_, err := Exec(user.UserID, args)
	return err
}
func (l *LocalUsers) InsertLoginUser(user *model_struct.LocalUser) error {
	_, err := Exec(utils.StructToJsonString(user))
	return err
}
