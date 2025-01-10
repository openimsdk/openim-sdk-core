//go:build !js
// +build !js

package db

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) initSuperLocalErrChatLog(ctx context.Context, groupID string) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	if !d.conn.WithContext(ctx).Migrator().HasTable(utils.GetErrTableName(groupID)) {
		d.conn.WithContext(ctx).Table(utils.GetErrTableName(groupID)).AutoMigrate(&model_struct.LocalErrChatLog{})
	}
}
func (d *DataBase) SuperBatchInsertExceptionMsg(ctx context.Context, MessageList []*model_struct.LocalErrChatLog, groupID string) error {
	if MessageList == nil {
		return nil
	}
	d.initSuperLocalErrChatLog(ctx, groupID)
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetConversationTableName(groupID)).Create(MessageList).Error, "BatchInsertMessageList failed")
}
func (d *DataBase) GetAbnormalMsgSeq(ctx context.Context) (int64, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var seq int64
	return seq, errs.WrapMsg(d.conn.WithContext(ctx).Model(model_struct.LocalErrChatLog{}).Select("IFNULL(max(seq),0)").Find(&seq).Error, "GetAbnormalMsgSeq")
}
func (d *DataBase) GetAbnormalMsgSeqList(ctx context.Context) ([]int64, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var seqList []int64
	return seqList, errs.WrapMsg(d.conn.WithContext(ctx).Model(model_struct.LocalErrChatLog{}).Select("seq").Find(&seqList).Error, "GetAbnormalMsgSeqList")
}
func (d *DataBase) BatchInsertExceptionMsg(ctx context.Context, messageList []*model_struct.LocalErrChatLog) error {
	if messageList == nil {
		return nil
	}
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return errs.WrapMsg(d.conn.WithContext(ctx).Create(messageList).Error, "BatchInsertMessageList failed")
}
func (d *DataBase) BatchInsertExceptionMsgController(ctx context.Context, messageList []*model_struct.LocalErrChatLog) error {
	if len(messageList) == 0 {
		return nil
	}
	switch messageList[len(messageList)-1].SessionType {
	case constant.ReadGroupChatType:
		return d.SuperBatchInsertExceptionMsg(ctx, messageList, messageList[len(messageList)-1].RecvID)
	default:
		return d.BatchInsertExceptionMsg(ctx, messageList)
	}
}
func (d *DataBase) GetConversationAbnormalMsgSeq(ctx context.Context, conversationID string) (int64, error) {
	var seq int64
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	if !d.conn.WithContext(ctx).Migrator().HasTable(utils.GetErrTableName(conversationID)) {
		return 0, nil
	}
	return seq, errs.WrapMsg(d.conn.WithContext(ctx).Table(utils.GetErrTableName(conversationID)).Select("IFNULL(max(seq),0)").Find(&seq).Error, "GetConversationNormalMsgSeq")
}
