package key

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/local_container"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/log"
)

type Key struct {
	//db          *db.DataBase
	db          db_interface.DataBase
	LoginUserId string
	ApiAddr     string
}

func NewKey(db db_interface.DataBase, loginUserId, ApiAddr string) *Key {
	return &Key{db: db, LoginUserId: loginUserId, ApiAddr: ApiAddr}
}

func (k Key) GetKey(sessionID string, sessionType int32) (*model_struct.LocalKey, error) {
	key, err := k.db.GetLocalKeyBySessionID(sessionID, sessionType)
	if err != nil {
		getKey, err2 := local_container.SdkGetKey(k.ApiAddr, k.LoginUserId, sessionID, sessionType)
		if err2 != nil {
			log.Error("key GetKey err " + err2.Error())
			return nil, err2
		}
		localKey := model_struct.LocalKey{
			SessionID:   sessionID,
			SessionKey:  getKey,
			SessionType: sessionType,
		}
		k.db.InsertLocalKey(&localKey)
		return &localKey, nil
	}
	return key, nil
}
func (k Key) GetAllKey() (*[]model_struct.LocalKey, error) {
	log.Info("login GetAllKey start************************")
	keys, err := k.db.GetAllLocalKey()
	if err != nil {
		log.Error("key GetAllKey err ", err.Error())
		return nil, err
	}
	return keys, nil
}

func (k Key) SynAllKey() {
	keysResp, err := local_container.SdkGetAllKey(k.ApiAddr, k.LoginUserId)
	if err != nil {
		log.Error("SynAllKey err LoginUserId", err.Error(), k.LoginUserId)
		return
	}
	log.NewInfo("SynAllKey start len ", len(keysResp.Keys))
	var localKeys []model_struct.LocalKey
	//copier.Copy(&localKeys, keysResp.Keys)
	for _, key := range keysResp.Keys {
		var localKey model_struct.LocalKey
		localKey.SessionID = key.SessionID
		localKey.SessionKey = key.SessionKey
		localKey.SessionType = key.SessionType
		localKeys = append(localKeys, localKey)
	}
	k.db.InsertAllLocalKey(&localKeys)
}

func (k Key) AddKey(sessionID, sessionKey string, sessionType int32) error {
	localKey := model_struct.LocalKey{}
	localKey.SessionKey = sessionID
	localKey.SessionKey = sessionKey
	localKey.SessionType = sessionType
	err := k.db.InsertLocalKey(&localKey)
	if err != nil {
		log.Error("AddKey err ", err.Error(), localKey)
		return err
	}
	return nil
}
