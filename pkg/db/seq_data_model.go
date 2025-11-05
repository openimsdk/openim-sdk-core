//go:build !js
// +build !js

package db

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) GetMinSeq(ctx context.Context, ID string) (uint32, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var seqData model_struct.LocalSeq
	return seqData.MinSeq, errs.WrapMsg(d.conn.WithContext(ctx).First(&seqData).Error, "GetMinSeq failed")
}

func (d *DataBase) SetMinSeq(ctx context.Context, ID string, minSeq uint32) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	seqData := model_struct.LocalSeq{ID: ID, MinSeq: minSeq}
	t := d.conn.WithContext(ctx).Updates(&seqData)
	if t.RowsAffected == 0 {
		return errs.WrapMsg(d.conn.WithContext(ctx).Create(seqData).Error, "Updates failed")
	} else {
		return errs.WrapMsg(t.Error, "SetMinSeq failed")
	}
}

func (d *DataBase) GetUserMinSeq(ctx context.Context) (uint32, error) {
	return 0, nil
}

func (d *DataBase) GetGroupMinSeq(ctx context.Context, groupID string) (uint32, error) {
	return 0, nil
}
