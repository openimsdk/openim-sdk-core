package db

import "open_im_sdk/pkg/utils"

func (d *DataBase) GetAbnormalMsgSeq() (uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seq uint32
	err := d.conn.Model(LocalErrChatLog{}).Select("max(seq)").Find(&seq).Error
	return seq, utils.Wrap(err, "GetAbnormalMsgSeq")
}
func (d *DataBase) BatchInsertExceptionMsgToErrorChatLog(MessageList []*LocalErrChatLog) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(MessageList).Error, "BatchInsertMessageList failed")
}
