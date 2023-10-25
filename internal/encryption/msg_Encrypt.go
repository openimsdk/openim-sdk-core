package encryption

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/aes"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/log"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

func msgEncryption(Message *sdk_struct.MsgStruct, key string) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("MsgEncryption err ", err)
		}
	}()
	switch Message.ContentType {
	//文本消息解密
	case constant.Text:
		msgStructTextEncryption(Message, key)
		//引用消息解密
	case constant.Quote:

		//@消息解密 AtElem
	case constant.AtText:

		//合并消息解密 MergeElem
	case constant.Merger:
		log.Info("MsgStructDecryption MsgStruct_mergerDecryption start******")
	}
}
func msgStructTextEncryption(message *sdk_struct.MsgStruct, key string) {
	byAes, err := aes.EncryptByAes([]byte(message.Content), []byte(key))
	if err != nil {
		log.Error("MsgStruct_textDecryption err ", key, message.Content)
	} else {
		message.Content = byAes
		message.Encryption = true
		message.EncryptionMode = "aes"
	}
}
