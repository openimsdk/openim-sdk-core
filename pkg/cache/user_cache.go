package cache

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/tools/utils/datautil"
)

func NewUserCache[K comparable, V any](
	getKeyFunc func(value V) K,
	batchDBFunc func(ctx context.Context, keys []K) ([]V, error),
	singleDBFunc func(ctx context.Context, keys K) (V, error),
	queryFunc func(ctx context.Context, keys []K) ([]V, error),
) *UserCache[K, V] {
	return &UserCache[K, V]{
		Cache:        NewCache[K, V](),
		getKeyFunc:   getKeyFunc,
		batchDBFunc:  batchDBFunc,
		singleDBFunc: singleDBFunc,
		queryFunc:    queryFunc,
	}
}

type UserCache[K comparable, V any] struct {
	*Cache[K, V]
	getKeyFunc   func(value V) K
	batchDBFunc  func(ctx context.Context, keys []K) ([]V, error)
	singleDBFunc func(ctx context.Context, keys K) (V, error)
	queryFunc    func(ctx context.Context, keys []K) ([]V, error)
}

func (m *UserCache[K, V]) BatchFetch(ctx context.Context, keys []K) (map[K]V, error) {
	var (
		res       = make(map[K]V)
		queryKeys []K
	)

	for _, key := range keys {
		if data, ok := m.Load(key); ok {
			res[key] = data
		} else {
			queryKeys = append(queryKeys, keys...)
		}
	}

	writeData, err := m.batchFetch(ctx, queryKeys)
	if err != nil {
		return nil, err
	}

	for i, data := range writeData {
		res[m.getKeyFunc(data)] = writeData[i]
		m.Store(m.getKeyFunc(data), writeData[i])
	}

	return res, nil
}

func (m *UserCache[K, V]) Fetch(ctx context.Context, key K) (V, error) {
	var nilData V

	if data, ok := m.Load(key); ok {
		return data, nil
	}

	fetchedData, err := m.fetch(ctx, key)
	if err != nil {
		return nilData, err
	}
	m.Store(key, fetchedData)
	return fetchedData, nil
}

func (m *UserCache[K, V]) batchFetch(ctx context.Context, keys []K) ([]V, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	var (
		queryKeys = keys
		writeData []V
	)

	if m.batchDBFunc != nil {
		dbData, err := m.batchDBFunc(ctx, queryKeys)
		if err != nil {
			return nil, err
		}
		writeData = dbData
		queryKeys = datautil.SliceSubAny(queryKeys, dbData, m.getKeyFunc)
	}

	if len(queryKeys) == 0 {
		return writeData, nil
	}

	if m.queryFunc != nil {
		queryData, err := m.queryFunc(ctx, queryKeys)
		if err != nil {
			return nil, err
		}
		if len(queryData) == 0 {
			return writeData, sdkerrs.ErrUserIDNotFound.WrapMsg("fetch data not found", "keys", keys)
		}
		writeData = append(writeData, queryData...)
	}

	return writeData, nil
}
func (m *UserCache[K, V]) fetch(ctx context.Context, key K) (V, error) {
	var writeData V
	if m.singleDBFunc != nil {
		dbData, err := m.singleDBFunc(ctx, key)
		if err == nil {
			return dbData, nil
		}
	}
	if m.queryFunc != nil {
		queryData, err := m.queryFunc(ctx, []K{key})
		if err != nil {
			return writeData, err
		}
		if len(queryData) > 0 {
			return queryData[0], nil
		}
		return writeData, sdkerrs.ErrUserIDNotFound.WrapMsg("fetch data not found", "key", key)
	}
	return writeData, nil
}
