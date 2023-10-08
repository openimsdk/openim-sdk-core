package decryption

import (
	"github.com/OpenIMSDK/tools/utils"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/aes"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/log"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

func MsgDecryption(Message *model_struct.LocalChatLog, key string) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("MsgDecryption err ", err)
		}
	}()
	switch Message.ContentType {
	//文本消息解密
	case constant.Text:
		textDecryption(Message, key)
		//引用消息解密
	case constant.Quote:
		quoteDecryption2(Message, key)
		//@消息解密 AtElem
	case constant.AtText:
		atTextDecryption(Message, key)
		//合并消息解密 MergeElem
	case constant.Merger:
		log.Info("MsgDecryption mergerDecryption start******")
		mergerDecryption2(Message, key)
	}
}
func textDecryption(Message *model_struct.LocalChatLog, key string) {
	byAes, err := aes.DecryptByAes(Message.Content, []byte(key))
	if err != nil {
		log.Error("textDecryption err ", key, Message.Content)
	} else {
		Message.Content = string(byAes)
	}
}
func quoteDecryption(Message *model_struct.LocalChatLog, key string) {
	quoteElem := sdk_struct.MsgStruct{}.QuoteElem
	var msgS sdk_struct.MsgStruct
	quoteElem.QuoteMessage = &msgS
	utils.JsonStringToStruct(Message.Content, &quoteElem)
	byAes_quote_text, err1 := aes.DecryptByAes(quoteElem.Text, []byte(key))
	byAes_quote_msg_cont, err2 := aes.DecryptByAes(quoteElem.QuoteMessage.Content, []byte(key))
	if err2 != nil && err1 != nil {
		log.Error("quoteDecryption  err ", key, quoteElem.QuoteMessage.Content)
	} else {
		quoteElem.Text = string(byAes_quote_text)
		quoteElem.QuoteMessage.Content = string(byAes_quote_msg_cont)
		Message.Content = utils.StructToJsonString(quoteElem)
	}
}
func quoteDecryption2(Message *model_struct.LocalChatLog, key string) {
	quoteElem := sdk_struct.MsgStruct{}.QuoteElem
	var msgS sdk_struct.MsgStruct
	quoteElem.QuoteMessage = &msgS
	utils.JsonStringToStruct(Message.Content, &quoteElem)
	byAes_quote_text, err1 := aes.DecryptByAes(quoteElem.Text, []byte(key))
	MsgStructDecryption(quoteElem.QuoteMessage, key)
	if err1 != nil {
		log.Error("quoteDecryption2  err ", key, quoteElem.QuoteMessage.Content)
		return
	} else {
		quoteElem.Text = string(byAes_quote_text)
		Message.Content = utils.StructToJsonString(quoteElem)
	}
}
func atTextDecryption(Message *model_struct.LocalChatLog, key string) {
	atElem := sdk_struct.MsgStruct{}.AtTextElem
	utils.JsonStringToStruct(Message.Content, &atElem)
	byAes_atElem_test, err := aes.DecryptByAes(atElem.Text, []byte(key))
	if err != nil {
		log.Error("atTextDecryption err ", key)
		return
	} else {
		atElem.Text = string(byAes_atElem_test)
		Message.Content = utils.StructToJsonString(atElem)
	}
	return
}
func mergerDecryption(Message *model_struct.LocalChatLog, key string) {
	mergerElem := sdk_struct.MsgStruct{}.MergeElem
	var structs []*sdk_struct.MsgStruct
	var entities []*sdk_struct.MessageEntity
	mergerElem.MultiMessage = structs
	mergerElem.MessageEntityList = entities
	utils.JsonStringToStruct(Message.Content, &mergerElem)
	for _, msgStruct := range mergerElem.MultiMessage {
		MsgStructDecryption(msgStruct, key)
	}
	Message.Content = utils.StructToJsonString(mergerElem)
}

func mergerDecryption2(Message *model_struct.LocalChatLog, key string) {
	msg := new(sdk_struct.MsgStruct)
	utils.JsonStringToStruct(Message.Content, &msg.MergeElem)
	for _, msgStruct := range msg.MergeElem.MultiMessage {
		MsgStructDecryption(msgStruct, key)
	}
	Message.Content = utils.StructToJsonString(msg.MergeElem)
}
