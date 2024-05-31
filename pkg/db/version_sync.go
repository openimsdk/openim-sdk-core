package db

import (
	"context"
	"errors"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"gorm.io/gorm"
)

func (d *DataBase) GetVersionSync(ctx context.Context, key string) (*model_struct.LocalVersionSync, error) {
	d.versionMtx.RLock()
	defer d.versionMtx.RUnlock()
	var res model_struct.LocalVersionSync
	return &res, utils.Wrap(d.conn.WithContext(ctx).Where("`key` = ? ", key).Take(&res).Error, "GetVersionSync failed")
}

func (d *DataBase) SetVersionSync(ctx context.Context, lv *model_struct.LocalVersionSync) error {
	d.versionMtx.Lock()
	defer d.versionMtx.Unlock()
	res := d.conn.WithContext(ctx).Where("`key` = ?", lv.Key).Save(lv)
	if res.Error != nil {
		return utils.Wrap(res.Error, "SetVersionSync failed")
	}
	if res.RowsAffected > 0 {
		return nil
	}
	if err := d.conn.WithContext(ctx).Create(lv).Error; err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		return utils.Wrap(err, "SetVersionSync failed")
	}
	return nil
}
