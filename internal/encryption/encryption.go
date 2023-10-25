package encryption

import (
	"github.com/openimsdk/openim-sdk-core/v3/internal/aes_key"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/log"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

type Encryption struct {
	IsEncryption   bool
	EncryptionMode string
	mode           EncModes
}

func (e Encryption) Mode() EncModes {
	return e.mode
}
func NewEncryption(isEncryption bool, mode string, key *aes_key.AesKey) *Encryption {
	if isEncryption {
		switch mode {
		case "aes":
			return &Encryption{isEncryption, mode, NewAesEncryption(key)}
		}
	}
	return &Encryption{isEncryption, mode, nil}
}

type EncModes interface {
	EncryptionMsg(msg *sdk_struct.MsgStruct)
}

type AesEncryption struct {
	key *aes_key.AesKey
}

func NewAesEncryption(key *aes_key.AesKey) EncModes {
	return &AesEncryption{key: key}
}
func (a *AesEncryption) EncryptionMsg(msg *sdk_struct.MsgStruct) {
	switch msg.ContentType {
	case constant.Text:
		key, err2 := a.key.GetKey(msg.SessionType, msg.GroupID, msg.SendID, msg.RecvID)
		if err2 != nil {
			log.Error("", "a.key.GetKey err ", err2.Error(), msg.SessionType, msg.GroupID, msg.SendID, msg.RecvID)
			return
		}
		//byAes, err := aes.EncryptByAes([]byte(msg.Content), []byte(key))
		//if err != nil {
		//	log.Error("", "aes.EncryptByAes err ", err.Error(), msg.ClientMsgID)
		//	return
		//}
		//msg.Content = byAes
		msgEncryption(msg, key)
	}
}
