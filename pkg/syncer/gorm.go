package syncer

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

type GormImpl struct {
	db *gorm.DB
}

func (g *GormImpl) where(ctx context.Context, where *Where) *gorm.DB {
	db := g.db.WithContext(ctx)
	for i := range where.Columns {
		db = db.Where(fmt.Sprintf("`%s` = ?", where.Columns[i].Name), where.Columns[i].Value)
	}
	return db
}

func (g *GormImpl) Delete(ctx context.Context, m reflect.Type, where []*Where) error {
	for _, w := range where {
		if err := g.where(ctx, w).Delete(reflect.New(m).Interface()).Error; err != nil {
			return err
		}
	}
	return nil
}

func (g *GormImpl) Update(ctx context.Context, m reflect.Type, where *Where, data map[string]any) error {
	return g.where(ctx, where).Model(reflect.New(m).Interface()).Updates(data).Error
}

func (g *GormImpl) FindOffset(ctx context.Context, m reflect.Type, where *Where, offset int, limit int) (any, error) {
	v := reflect.New(reflect.SliceOf(m)).Interface()
	return v, g.where(ctx, where).Offset(offset).Limit(limit).Find(v).Error
}

func (g *GormImpl) Insert(ctx context.Context, ms any) error {
	return g.db.WithContext(ctx).Create(ms).Error
}
