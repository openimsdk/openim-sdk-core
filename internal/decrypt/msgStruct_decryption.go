package decrypt

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/aes"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/log"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

func msgStructDecryption(Message *sdk_struct.MsgStruct, key string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error("MsgDecryption err ", err)
		}
	}()
	switch Message.ContentType {
	//文本消息解密
	case constant.Text:
		return msgStructTextDecryption(Message, key)
	//引用消息解密
	case constant.Quote:
		//@消息解密 AtElem
	case constant.AtText:
		//合并消息解密 MergeElem
	case constant.Merger:
		log.Info("MsgStructDecryption MsgStruct_mergerDecryption start******")
	}
	return nil
}
func msgStructTextDecryption(Message *sdk_struct.MsgStruct, key string) error {
	byAes, err := aes.DecryptByAes(Message.Content, []byte(key))
	if err != nil {
		log.Error("MsgStruct_textDecryption err ", key, Message.Content)
		return err
	} else {
		Message.Content = string(byAes)
		Message.Encryption = false
		return nil
	}
}
