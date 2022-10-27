package indexdb

import (
	"errors"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

var ErrType = errors.New("from javascript data type err")
var PrimaryKeyNull = errors.New("primary key is null err")

type LocalUsers struct {
}

func (l *LocalUsers) GetLoginUser(userID string) (*model_struct.LocalUser, error) {
	user, err := Exec(userID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := user.(string); ok {
			result := model_struct.LocalUser{}
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

func (l *LocalUsers) UpdateLoginUser(user *model_struct.LocalUser) error {
	_, err := Exec(user)
	return err

}
func (l *LocalUsers) UpdateLoginUserByMap(user *model_struct.LocalUser, args map[string]interface{}) error {
	_, err := Exec(user.UserID, args)
	return err
}
func (l *LocalUsers) InsertLoginUser(user *model_struct.LocalUser) error {
	_, err := Exec(utils.StructToJsonString(user))
	return err
}
