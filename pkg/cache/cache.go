package cache

import "sync"

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{}
}

// Cache is a Generic sync.Map structure.
type Cache[K comparable, V any] struct {
	m sync.Map
}

// Load returns the value stored in the map for a key, or nil if no value is present.
func (c *Cache[K, V]) Load(key K) (value V, ok bool) {
	rawValue, ok := c.m.Load(key)
	if !ok {
		return
	}
	return rawValue.(V), ok
}

// Store sets the value for a key.
func (c *Cache[K, V]) Store(key K, value V) {
	c.m.Store(key, value)
}

// StoreWithFunc stores the value for a key only if the provided condition function returns true.
func (c *Cache[K, V]) StoreWithFunc(key K, value V, condition func(key K, value V) bool) {
	if condition(key, value) {
		c.m.Store(key, value)
	}
}

// StoreAll sets all value by f's key.
func (c *Cache[K, V]) StoreAll(f func(value V) K, values []V) {
	for _, v := range values {
		c.m.Store(f(v), v)
	}
}

// LoadOrStore returns the existing value for the key if present.
func (c *Cache[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	rawValue, loaded := c.m.LoadOrStore(key, value)
	return rawValue.(V), loaded
}

// Delete deletes the value for a key.
func (c *Cache[K, V]) Delete(key K) {
	c.m.Delete(key)
}

// DeleteAll deletes all values.
func (c *Cache[K, V]) DeleteAll() {
	c.m.Range(func(key, value interface{}) bool {
		c.m.Delete(key)
		return true
	})
}

// DeleteCon deletes the value for a key only if the provided condition function returns true.
func (c *Cache[K, V]) DeleteCon(condition func(key K, value V) bool) {
	c.m.Range(func(rawKey, rawValue interface{}) bool {
		if condition(rawKey.(K), rawValue.(V)) {
			c.m.Delete(rawKey)
		}
		return true // Continue iteration
	})
}

// RangeAll returns all values in the map.
func (c *Cache[K, V]) RangeAll() (values []V) {
	c.m.Range(func(rawKey, rawValue interface{}) bool {
		values = append(values, rawValue.(V))
		return true
	})
	return values
}

// RangeCon returns values in the map that satisfy condition f.
func (c *Cache[K, V]) RangeCon(f func(key K, value V) bool) (values []V) {
	c.m.Range(func(rawKey, rawValue interface{}) bool {
		if f(rawKey.(K), rawValue.(V)) {
			values = append(values, rawValue.(V))
		}
		return true
	})
	return values
}
