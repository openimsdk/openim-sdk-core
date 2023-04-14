//go:build !js
// +build !js

package db

import (
	"context"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) initSuperLocalErrChatLog(ctx context.Context, groupID string) {
	if !d.conn.WithContext(ctx).Migrator().HasTable(utils.GetErrSuperGroupTableName(groupID)) {
		d.conn.WithContext(ctx).Table(utils.GetErrSuperGroupTableName(groupID)).AutoMigrate(&model_struct.LocalErrChatLog{})
	}
}
func (d *DataBase) SuperBatchInsertExceptionMsg(ctx context.Context, MessageList []*model_struct.LocalErrChatLog, groupID string) error {
	if MessageList == nil {
		return nil
	}
	d.initSuperLocalErrChatLog(ctx, groupID)
	return utils.Wrap(d.conn.WithContext(ctx).Table(utils.GetSuperGroupTableName(groupID)).Create(MessageList).Error, "BatchInsertMessageList failed")
}
func (d *DataBase) GetAbnormalMsgSeq(ctx context.Context) (uint32, error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var seq uint32
	err := d.conn.WithContext(ctx).Model(model_struct.LocalErrChatLog{}).Select("IFNULL(max(seq),0)").Find(&seq).Error
	return seq, utils.Wrap(err, "GetAbnormalMsgSeq")
}
func (d *DataBase) GetAbnormalMsgSeqList(ctx context.Context) ([]uint32, error) {
	var seqList []uint32
	err := d.conn.WithContext(ctx).Model(model_struct.LocalErrChatLog{}).Select("seq").Find(&seqList).Error
	return seqList, utils.Wrap(err, "GetAbnormalMsgSeqList")
}
func (d *DataBase) BatchInsertExceptionMsg(ctx context.Context, messageList []*model_struct.LocalErrChatLog) error {
	if messageList == nil {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Create(MessageList).Error, "BatchInsertMessageList failed")
}
func (d *DataBase) BatchInsertExceptionMsgController(ctx context.Context, messageList []*model_struct.LocalErrChatLog) error {
	if len(messageList) == 0 {
		return nil
	}
	switch messageList[len(messageList)-1].SessionType {
	case constant.SuperGroupChatType:
		return d.SuperBatchInsertExceptionMsg(ctx, MessageList, MessageList[len(MessageList)-1].RecvID)
	default:
		return d.BatchInsertExceptionMsg(ctx, MessageList)
	}
}
func (d *DataBase) GetSuperGroupAbnormalMsgSeq(ctx context.Context, groupID string) (uint32, error) {
	var seq uint32
	if !d.conn.WithContext(ctx).Migrator().HasTable(utils.GetErrSuperGroupTableName(groupID)) {
		return 0, nil
	}
	err := d.conn.WithContext(ctx).Table(utils.GetErrSuperGroupTableName(groupID)).Select("IFNULL(max(seq),0)").Find(&seq).Error
	return seq, utils.Wrap(err, "GetSuperGroupNormalMsgSeq")
}
