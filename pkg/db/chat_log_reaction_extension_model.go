package db

import (
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetMessageReactionExtension(msgID string) (result *model_struct.LocalChatLogReactionExtensions, err error) {
	var l model_struct.LocalChatLogReactionExtensions
	return &l, utils.Wrap(d.conn.Where("client_msg_id = ?",
		msgID).Take(&l).Error, "GetMessageReactionExtension failed")
}

func (d *DataBase) InsertMessageReactionExtension(messageReactionExtension *model_struct.LocalChatLogReactionExtensions) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(messageReactionExtension).Error, "InsertMessageReactionExtension failed")
}
