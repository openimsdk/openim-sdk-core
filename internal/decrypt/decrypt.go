package decrypt

import (
	"github.com/openimsdk/openim-sdk-core/v3/internal/aes_key"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/log"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

type Decrypt struct {
	key *aes_key.AesKey
}

func NewDecrypt(key *aes_key.AesKey) *Decrypt {
	return &Decrypt{key: key}
}
func (d *Decrypt) DecryptMsg(msg *sdk_struct.MsgStruct) {
	switch msg.EncryptionMode {
	case "aes":
		key, err2 := d.key.GetKey(msg.SessionType, msg.GroupID, msg.SendID, msg.RecvID)
		if err2 != nil {
			log.Error("", "a.key.GetKey err ", err2.Error(), msg.SessionType, msg.GroupID, msg.SendID, msg.RecvID)
			return
		}
		AesDecrypt(msg, key)
	}
}

func AesDecrypt(msg *sdk_struct.MsgStruct, key string) {
	switch msg.ContentType {
	case constant.Text:
		msgStructDecryption(msg, key)
	}
}
