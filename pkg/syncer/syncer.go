package syncer

import (
	"context"
	"reflect"
)

func New[T any, V comparable](
	insert func(ctx context.Context, value T) error,
	delete func(ctx context.Context, value T) error,
	update func(ctx context.Context, server T, local T) error,
	uuid func(value T) V,
	equal func(a T, b T) bool,
	opt *Option[T],
) *Syncer[T, V] {
	if insert == nil || update == nil || delete == nil || uuid == nil {
		panic("invalid params")
	}
	return &Syncer[T, V]{
		insert: insert,
		update: update,
		delete: delete,
		uuid:   uuid,
		equal:  equal,
		opt:    opt,
	}
}

type Syncer[T any, V comparable] struct {
	insert func(ctx context.Context, value T) error
	update func(ctx context.Context, server T, local T) error
	delete func(ctx context.Context, value T) error
	equal  func(server T, local T) bool
	uuid   func(value T) V
	opt    *Option[T]
}

func (s *Syncer[T, V]) eq(server T, local T) bool {
	if s.equal != nil {
		return s.equal(server, local)
	}
	return reflect.DeepEqual(server, local)
}

func (s *Syncer[T, V]) onNotice(ctx context.Context, state int, server, local T, opt *Option[T]) error {
	if s.opt != nil {
		if err := s.opt.on(ctx, state, server, local); err != nil {
			return err
		}
	}
	if opt != nil {
		if err := opt.on(ctx, state, server, local); err != nil {
			return err
		}
	}
	return nil
}

func (s *Syncer[T, V]) Sync(ctx context.Context, serverData []T, localData []T, notice *Option[T]) error {
	if len(serverData) == 0 && len(localData) == 0 {
		return nil
	}
	localMap := make(map[V]T)
	for i, item := range localData {
		localMap[s.uuid(item)] = localData[i]
	}
	for i := range serverData {
		server := serverData[i]
		id := s.uuid(server)
		local, ok := localMap[id]
		if !ok {
			if err := s.insert(ctx, server); err != nil {
				return err
			}
			if err := s.onNotice(ctx, Insert, server, local, notice); err != nil {
				return err
			}
			continue
		}
		delete(localMap, id)
		if s.eq(server, local) {
			if err := s.onNotice(ctx, Unchanged, server, local, notice); err != nil {
				return err
			}
			continue
		}
		if err := s.update(ctx, server, local); err != nil {
			return err
		}
		if err := s.onNotice(ctx, Update, server, local, notice); err != nil {
			return err
		}
	}
	for id := range localMap {
		local := localMap[id]
		if err := s.delete(ctx, local); err != nil {
			return err
		}
		var server T
		if err := s.onNotice(ctx, Delete, server, local, notice); err != nil {
			return err
		}
	}
	return nil
}

func (s *Syncer[T, V]) NewOption() *Option[T] {
	return &Option[T]{}
}
