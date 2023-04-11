package syncer

import "reflect"

type DBInterface interface {
	// Insert 插入数据 ms 是model层的结构体切片
	Insert(ms any) error // []Struct
	// Delete reflect.Type 返回的是最底层的结构体类型
	Delete(m reflect.Type, where []*Where) error
	// Update reflect.Type 返回的是最底层的结构体类型
	Update(m reflect.Type, where *Where, data map[string]any) error
	// FindOffset reflect.Type 返回的是最底层的结构体类型
	FindOffset(m reflect.Type, where *Where, offset int, limit int) (any, error) // []Struct
}
