package cache

import (
	"context"
	"github.com/openimsdk/tools/utils/datautil"
)

func NewManager[K comparable, V any](
	getKeyFunc func(value V) K,
	dbFunc func(ctx context.Context, keys []K) ([]V, error),
	queryFunc func(ctx context.Context, keys []K) ([]V, error),
) *Manager[K, V] {
	return &Manager[K, V]{
		Cache:      Cache[K, V]{},
		getKeyFunc: getKeyFunc,
		dbFunc:     dbFunc,
		queryFunc:  queryFunc,
	}
}

type Manager[K comparable, V any] struct {
	Cache[K, V]
	getKeyFunc func(value V) K
	dbFunc     func(ctx context.Context, keys []K) ([]V, error)
	queryFunc  func(ctx context.Context, keys []K) ([]V, error)
}

func (m *Manager[K, V]) MultiFetchGet(ctx context.Context, keys []K) (map[K]V, error) {
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

	writeData, err := m.Fetch(ctx, queryKeys)
	if err != nil {
		return nil, err
	}

	for i, data := range writeData {
		res[m.getKeyFunc(data)] = writeData[i]
		m.Store(m.getKeyFunc(data), writeData[i])
	}

	return res, nil
}

func (m *Manager[K, V]) FetchGet(ctx context.Context, key K) (V, error) {
	var nilData V

	if data, ok := m.Load(key); ok {
		return data, nil
	}

	fetchedData, err := m.Fetch(ctx, []K{key})
	if err != nil {
		return nilData, err
	}
	if len(fetchedData) > 0 {
		m.Store(key, fetchedData[0])
		return fetchedData[0], nil
	}

	// todo: return error or nilData?
	return nilData, nil
}

func (m *Manager[K, V]) Fetch(ctx context.Context, keys []K) ([]V, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	var (
		queryKeys = keys
		writeData []V
	)

	if m.dbFunc != nil {
		dbData, err := m.dbFunc(ctx, queryKeys)
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
		writeData = append(writeData, queryData...)
	}

	return writeData, nil
}
