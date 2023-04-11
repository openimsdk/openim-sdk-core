package syncer

import (
	"context"
	"reflect"
)

type DBInterface interface {
	// Insert 插入数据 ms 是model层的结构体切片
	Insert(ctx context.Context, ms any) error // []Struct
	// Delete reflect.Type 返回的是最底层的结构体类型
	Delete(ctx context.Context, m reflect.Type, where []*Where) error
	// Update reflect.Type 返回的是最底层的结构体类型
	Update(ctx context.Context, m reflect.Type, where *Where, data map[string]any) error
	// FindOffset reflect.Type 返回的是最底层的结构体类型
	FindOffset(ctx context.Context, m reflect.Type, where *Where, offset int, limit int) (any, error) // []Struct
}
