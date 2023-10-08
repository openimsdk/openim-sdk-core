package decryption

import (
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/aes"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/log"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

func MsgStructDecryption(Message *sdk_struct.MsgStruct, key string) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("MsgDecryption err ", err)
		}
	}()
	switch Message.ContentType {
	//文本消息解密
	case constant.Text:
		MsgStruct_textDecryption(Message, key)
		//引用消息解密
	case constant.Quote:
		MsgStruct_quoteDecryption(Message, key)
		//@消息解密 AtElem
	case constant.AtText:
		MsgStruct_atTextDecryption(Message, key)
		//合并消息解密 MergeElem
	case constant.Merger:
		log.Info("MsgStructDecryption MsgStruct_mergerDecryption start******")
		MsgStruct_mergerDecryption(Message, key)
	}
}
func MsgStruct_textDecryption(Message *sdk_struct.MsgStruct, key string) {
	byAes, err := aes.DecryptByAes(Message.TextElem.Content, []byte(key))
	if err != nil {
		log.Error("MsgStruct_textDecryption err ", key, Message.TextElem.Content)
	} else {
		Message.TextElem.Content = string(byAes)
	}
}
func MsgStruct_quoteDecryption(Message *sdk_struct.MsgStruct, key string) {
	quoteElem := sdk_struct.MsgStruct{}.QuoteElem
	var msgS sdk_struct.MsgStruct
	quoteElem.QuoteMessage = &msgS
	utils.JsonStringToStruct(Message.Content, &quoteElem)
	byAes_quote_text, err1 := aes.DecryptByAes(quoteElem.Text, []byte(key))
	if quoteElem.QuoteMessage != nil && quoteElem.QuoteMessage.Content != "" {
		MsgStructDecryption(quoteElem.QuoteMessage, key)
	}
	if err1 != nil {
		log.Error("MsgStruct_quoteDecryption  err ", key, quoteElem.QuoteMessage.Content)
	} else {
		quoteElem.Text = string(byAes_quote_text)
		Message.Content = utils.StructToJsonString(quoteElem)
	}
}
func MsgStruct_quoteDecryption2(Message *sdk_struct.MsgStruct, key string) {
	if &Message.QuoteElem != nil {
		byAes, err := aes.DecryptByAes(Message.QuoteElem.Text, []byte(key))
		if err != nil {
			return
		}
		Message.QuoteElem.Text = string(byAes)
		MsgStructDecryption(Message.QuoteElem.QuoteMessage, key)
	}

	quoteElem := sdk_struct.MsgStruct{}.QuoteElem
	var msgS sdk_struct.MsgStruct
	quoteElem.QuoteMessage = &msgS
	utils.JsonStringToStruct(Message.Content, &quoteElem)
	byAes_quote_text, err1 := aes.DecryptByAes(quoteElem.Text, []byte(key))
	MsgStructDecryption(quoteElem.QuoteMessage, key)
	if err1 != nil {
		log.Error("MsgStruct_quoteDecryption  err ", key, quoteElem.QuoteMessage.Content)
		return
	} else {
		quoteElem.Text = string(byAes_quote_text)
		Message.Content = utils.StructToJsonString(quoteElem)
	}
}
func MsgStruct_atTextDecryption(Message *sdk_struct.MsgStruct, key string) {
	if &Message.AtTextElem != nil {
		byAes, err := aes.DecryptByAes(Message.AtTextElem.Text, []byte(key))
		if err != nil {
			return
		}
		Message.AtTextElem.Text = string(byAes)
	}

	atElem := sdk_struct.MsgStruct{}.AtTextElem
	utils.JsonStringToStruct(Message.Content, &atElem)
	byAes_atElem_test, err := aes.DecryptByAes(atElem.Text, []byte(key))
	if err != nil {
		log.Error("MsgStruct_atTextDecryption err ", key)
		return
	} else {
		atElem.Text = string(byAes_atElem_test)
		Message.Content = utils.StructToJsonString(atElem)
	}
}
func MsgStruct_mergerDecryption(Message *sdk_struct.MsgStruct, key string) {
	msg := new(sdk_struct.MsgStruct)
	utils.JsonStringToStruct(Message.Content, &msg.MergeElem)
	for _, msgStruct := range msg.MergeElem.MultiMessage {
		MsgStructDecryption(msgStruct, key)
	}
	Message.Content = utils.StructToJsonString(msg.MergeElem)
}
