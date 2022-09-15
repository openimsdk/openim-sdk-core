package indexdb

import (
	"errors"
	"open_im_sdk/pkg/db/model_struct"
)

type LocalUsers struct {
}

func (l *LocalUsers) GetLoginUser(userID string) (*model_struct.LocalUser, error) {
	user, err := Exec(userID)
	if err != nil {
		return nil, err
	} else {
		v, ok := user.(model_struct.LocalUser)
		if ok {
			return &v, nil
		} else {
			return nil, errors.New("type err")
		}
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
	_, err := Exec(user)
	return err
}
