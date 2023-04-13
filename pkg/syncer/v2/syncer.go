package v2

import "context"

//type Interface[T any, V comparable] interface {
//	OnInsert(ctx context.Context, value T) error
//	OnDelete(ctx context.Context, value T) error
//	OnUpdate(ctx context.Context, server T, local T) error
//	GetID(value T) V
//}
//
//type Equal[T any] interface {
//	Equal(a, b T) bool
//}
//
//type UnchangedInterface[T any] interface {
//	Unchanged(ctx context.Context, value T) error
//}
//
//type InsertInterface[T any] interface {
//	Insert(ctx context.Context, value T) error
//}
//
//type DeleteInterface[T any] interface {
//	Delete(ctx context.Context, value T) error
//}
//
//type UpdateInterface[T any] interface {
//	Update(ctx context.Context, server T, local T) error
//}

type SyncerInterface[T any, V comparable] interface {
	OnInsert(ctx context.Context, value T) error
	OnDelete(ctx context.Context, value T) error
	OnUpdate(ctx context.Context, server, local T) error
	GetID(value T) V
	Equal(a, b T) bool
	Unchanged(ctx context.Context, value T) error
	Insert(ctx context.Context, value T) error
	Delete(ctx context.Context, value T) error
	Update(ctx context.Context, server T, local T) error
	Sync(ctx context.Context, server, local []T) error
}

func Sync[T any, V comparable](ctx context.Context, serverData []T, localData []T, syncer SyncerInterface[T, V]) error {

	return nil
}
