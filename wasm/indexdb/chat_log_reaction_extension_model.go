//go:build js && wasm
// +build js,wasm

package indexdb

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

type LocalChatLogReactionExtensions struct {
}

func NewLocalChatLogReactionExtensions() *LocalChatLogReactionExtensions {
	return &LocalChatLogReactionExtensions{}
}

func (i *LocalChatLogReactionExtensions) GetMessageReactionExtension(ctx context.Context, clientMsgID string) (result *model_struct.LocalChatLogReactionExtensions, err error) {
	msg, err := Exec(clientMsgID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := msg.(string); ok {
			result := model_struct.LocalChatLogReactionExtensions{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, ErrType
		}
	}
}

func (i *LocalChatLogReactionExtensions) InsertMessageReactionExtension(ctx context.Context, messageReactionExtension *model_struct.LocalChatLogReactionExtensions) error {
	_, err := Exec(utils.StructToJsonString(messageReactionExtension))
	return err
}
func (i *LocalChatLogReactionExtensions) GetAndUpdateMessageReactionExtension(ctx context.Context, clientMsgID string, m map[string]*sdkws.KeyValue) error {
	_, err := Exec(clientMsgID, utils.StructToJsonString(m))
	return err
}
func (i *LocalChatLogReactionExtensions) DeleteAndUpdateMessageReactionExtension(ctx context.Context, clientMsgID string, m map[string]*sdkws.KeyValue) error {
	_, err := Exec(clientMsgID, utils.StructToJsonString(m))
	return err
}
func (i *LocalChatLogReactionExtensions) GetMultipleMessageReactionExtension(ctx context.Context, msgIDList []string) (result []*model_struct.LocalChatLogReactionExtensions, err error) {
	msgReactionExtensionList, err := Exec(utils.StructToJsonString(msgIDList))
	if err != nil {
		return nil, err
	} else {
		if v, ok := msgReactionExtensionList.(string); ok {
			var temp []model_struct.LocalChatLogReactionExtensions
			err := utils.JsonStringToStruct(v, &temp)
			if err != nil {
				return nil, err
			}
			for _, v := range temp {
				v1 := v
				result = append(result, &v1)
			}
			return result, err
		} else {
			return nil, ErrType
		}
	}
}
func (i *LocalChatLogReactionExtensions) DeleteMessageReactionExtension(ctx context.Context, msgID string) error {
	_, err := Exec(msgID)
	return err
}
func (i *LocalChatLogReactionExtensions) UpdateMessageReactionExtension(ctx context.Context, c *model_struct.LocalChatLogReactionExtensions) error {
	_, err := Exec(c.ClientMsgID, utils.StructToJsonString(c))
	return err
}
