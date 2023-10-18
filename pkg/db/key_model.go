package db

import (
	"errors"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

func (d *DataBase) InsertLocalKey(localKey *model_struct.LocalKey) error {

	if localKey.SessionKey == "" {
		return errors.New("SessionKey is null")
	}
	return utils.Wrap(d.conn.Create(localKey).Error, "db InsertLocalKey failed")
}
func (d *DataBase) InsertAllLocalKey(localKey *[]model_struct.LocalKey) error {

	return utils.Wrap(d.conn.Create(localKey).Error, "db InsertAllLocalKey failed")
}
func (d *DataBase) GetLocalKeyBySessionID(sessionID string, sessionType int32) (*model_struct.LocalKey, error) {

	var key model_struct.LocalKey
	return &key, utils.Wrap(d.conn.Where("session_id = ? and session_type = ?", sessionID, sessionType).Take(&key).Error, "db GetLocalKeyBySessionID failed")
}

func (d *DataBase) GetAllLocalKey() (*[]model_struct.LocalKey, error) {
	var key []model_struct.LocalKey
	return &key, utils.Wrap(d.conn.Find(&key).Error, "db GetAllLocalKey failed")
}
