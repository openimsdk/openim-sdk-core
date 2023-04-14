//go:build !js
// +build !js

package db

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetMinSeq(ctx context.Context, ID string) (uint32, error) {
	var seqData model_struct.LocalSeq
	return seqData.MinSeq, utils.Wrap(d.conn.WithContext(ctx).First(&seqData).Error, "GetMinSeq failed")
}

func (d *DataBase) SetMinSeq(ctx context.Context, ID string, minSeq uint32) error {
	seqData := model_struct.LocalSeq{ID: ID, MinSeq: minSeq}
	t := d.conn.WithContext(ctx).Updates(&seqData)
	if t.RowsAffected == 0 {
		return utils.Wrap(d.conn.WithContext(ctx).Create(seqData).Error, "Updates failed")
	} else {
		return utils.Wrap(t.Error, "SetMinSeq failed")
	}
}

func (d *DataBase) GetUserMinSeq(ctx context.Context) (uint32, error) {
	return 0, nil
}

func (d *DataBase) GetGroupMinSeq(ctx context.Context, groupID string) (uint32, error) {
	return 0, nil
}
