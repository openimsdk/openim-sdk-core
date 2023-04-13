package syncer

import "context"

type Option[T any] struct {
	unchanged []func(ctx context.Context, value T) error
	insert    []func(ctx context.Context, value T) error
	update    []func(ctx context.Context, server, local T) error
	delete    []func(ctx context.Context, value T) error
	change    []func(ctx context.Context, state int, server, local T) error
}

func (opt *Option[T]) Insert(fn func(ctx context.Context, value T) error) *Option[T] {
	opt.insert = append(opt.insert, fn)
	return opt
}

func (opt *Option[T]) Update(fn func(ctx context.Context, newValue, oldValue T) error) *Option[T] {
	opt.update = append(opt.update, fn)
	return opt
}

func (opt *Option[T]) Delete(fn func(ctx context.Context, value T) error) *Option[T] {
	opt.delete = append(opt.delete, fn)
	return opt
}

func (opt *Option[T]) Unchanged(fn func(ctx context.Context, value T) error) *Option[T] {
	opt.delete = append(opt.delete, fn)
	return opt
}

func (opt *Option[T]) Change(fn func(ctx context.Context, state int, newValue, oldValue T) error) *Option[T] {
	opt.change = append(opt.change, fn)
	return opt
}

func (opt *Option[T]) on(ctx context.Context, state int, server, local T) error {
	switch state {
	case Unchanged:
		for _, fn := range opt.unchanged {
			if err := fn(ctx, server); err != nil {
				return err
			}
		}
	case Insert:
		for _, fn := range opt.insert {
			if err := fn(ctx, server); err != nil {
				return err
			}
		}
	case Delete:
		for _, fn := range opt.delete {
			if err := fn(ctx, local); err != nil {
				return err
			}
		}
	case Update:
		for _, fn := range opt.update {
			if err := fn(ctx, server, local); err != nil {
				return err
			}
		}
	}
	if state != Unchanged {
		for _, fn := range opt.change {
			if err := fn(ctx, state, server, local); err != nil {
				return err
			}
		}
	}
	return nil
}
