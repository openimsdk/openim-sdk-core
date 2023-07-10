package db

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) initSuperLocalErrChatLog(groupID string) {
	if !d.conn.Migrator().HasTable(utils.GetErrSuperGroupTableName(groupID)) {
		d.conn.Table(utils.GetErrSuperGroupTableName(groupID)).AutoMigrate(&model_struct.LocalErrChatLog{})
	}
}
func (d *DataBase) SuperBatchInsertExceptionMsg(MessageList []*model_struct.LocalErrChatLog, groupID string) error {
	if MessageList == nil {
		return nil
	}
	d.initSuperLocalErrChatLog(groupID)
	return utils.Wrap(d.conn.Table(utils.GetSuperGroupTableName(groupID)).Create(MessageList).Error, "BatchInsertMessageList failed")
}
func (d *DataBase) GetAbnormalMsgSeq() (uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seq uint32
	err := d.conn.Model(model_struct.LocalErrChatLog{}).Select("IFNULL(max(seq),0)").Find(&seq).Error
	return seq, utils.Wrap(err, "GetAbnormalMsgSeq")
}
func (d *DataBase) GetAbnormalMsgSeqList() ([]uint32, error) {
	var seqList []uint32
	err := d.conn.Model(model_struct.LocalErrChatLog{}).Select("seq").Find(&seqList).Error
	return seqList, utils.Wrap(err, "GetAbnormalMsgSeqList")
}
func (d *DataBase) BatchInsertExceptionMsg(MessageList []*model_struct.LocalErrChatLog) error {
	if MessageList == nil {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.Create(MessageList).Error, "BatchInsertMessageList failed")
}
func (d *DataBase) BatchInsertExceptionMsgController(MessageList []*model_struct.LocalErrChatLog) error {
	if len(MessageList) == 0 {
		return nil
	}
	switch MessageList[len(MessageList)-1].SessionType {
	case constant.SuperGroupChatType:
		return d.SuperBatchInsertExceptionMsg(MessageList, MessageList[len(MessageList)-1].RecvID)
	default:
		return d.BatchInsertExceptionMsg(MessageList)
	}
}
func (d *DataBase) GetSuperGroupAbnormalMsgSeq(groupID string) (uint32, error) {
	var seq uint32
	if !d.conn.Migrator().HasTable(utils.GetErrSuperGroupTableName(groupID)) {
		return 0, nil
	}
	err := d.conn.Table(utils.GetErrSuperGroupTableName(groupID)).Select("IFNULL(max(seq),0)").Find(&seq).Error
	return seq, utils.Wrap(err, "GetSuperGroupNormalMsgSeq")
}
