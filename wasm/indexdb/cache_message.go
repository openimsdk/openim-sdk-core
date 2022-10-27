package indexdb

import "open_im_sdk/pkg/db/model_struct"

type LocalCacheMessage struct {
}

func (i *LocalCacheMessage) BatchInsertTempCacheMessageList(MessageList []*model_struct.TempCacheLocalChatLog) error {
	panic("implement me")
}

func (i *LocalCacheMessage) InsertTempCacheMessage(Message *model_struct.TempCacheLocalChatLog) error {
	panic("implement me")
}
