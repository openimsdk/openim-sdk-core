//go:build js && wasm
// +build js,wasm

package indexdb

import (
	"errors"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/indexdb/temp_struct"
	"time"
)

var ErrType = errors.New("from javascript data type err")
var PrimaryKeyNull = errors.New("primary key is null err")

type LocalUsers struct {
}

func NewLocalUsers() *LocalUsers {
	return &LocalUsers{}
}

func (l *LocalUsers) GetLoginUser(userID string) (*model_struct.LocalUser, error) {
	user, err := Exec(userID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := user.(string); ok {
			result := model_struct.LocalUser{}
			temp := temp_struct.LocalUser{}
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			result.UserID = temp.UserID
			result.Nickname = temp.Nickname
			result.FaceURL = temp.FaceURL
			result.Gender = temp.Gender
			result.PhoneNumber = temp.PhoneNumber
			result.Birth = temp.Birth
			result.Email = temp.Email
			result.CreateTime = temp.CreateTime
			result.AppMangerLevel = temp.AppMangerLevel
			result.Ex = temp.Ex
			result.AttachedInfo = temp.Ex
			result.GlobalRecvMsgOpt = temp.GlobalRecvMsgOpt
			time, err := utils.TimeStringToTime(temp.BirthTime)
			if err != nil {
				return nil, err
			}
			result.BirthTime = time
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
	if v, ok := args["birth_time"]; ok {
		if t, ok := v.(time.Time); ok {
			args["birth_time"] = utils.TimeToString(t)
		} else {
			return ErrType
		}
	}
	_, err := Exec(user.UserID, args)
	return err
}
func (l *LocalUsers) InsertLoginUser(user *model_struct.LocalUser) error {
	temp := temp_struct.LocalUser{}
	temp.UserID = user.UserID
	temp.Nickname = user.Nickname
	temp.FaceURL = user.FaceURL
	temp.Gender = user.Gender
	temp.PhoneNumber = user.PhoneNumber
	temp.Birth = user.Birth
	temp.Email = user.Email
	temp.CreateTime = user.CreateTime
	temp.AppMangerLevel = user.AppMangerLevel
	temp.Ex = user.Ex
	temp.AttachedInfo = user.Ex
	temp.GlobalRecvMsgOpt = user.GlobalRecvMsgOpt
	t := utils.TimeToString(user.BirthTime)
	temp.BirthTime = t
	_, err := Exec(utils.StructToJsonString(temp))
	return err
}
