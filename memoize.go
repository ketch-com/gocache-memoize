package memoize

import (
	"context"
	"errors"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	"golang.org/x/sync/singleflight"
)

// Memoizer provides a mechanism to memoize results of functions
type Memoizer[T any] struct {
	*cache.Cache[T]
	group singleflight.Group
}

// NewMemoizer returns a new Memoizer
func NewMemoizer[T any](cache *cache.Cache[T]) *Memoizer[T] {
	return &Memoizer[T]{
		Cache: cache,
	}
}

// Memoize returns the last value of the function fn when called with the given key. If not available in the
// cache, then the function is called and its result is cached.
func (m *Memoizer[T]) Memoize(ctx context.Context, key string, fn func(context.Context) (T, error)) (T, error) {
	value, err := m.Get(ctx, key)
	if err == nil {
		return value, nil
	} else if !errors.Is(err, store.NotFound{}) {
		return *new(T), err
	}

	data, err, _ := m.group.Do(key, func() (any, error) {
		v, err := fn(ctx)
		if err != nil {
			return nil, err
		}

		if err := m.Set(ctx, key, v); err != nil {
			return nil, err
		}

		return v, nil
	})
	if err != nil {
		return *new(T), err
	}

	if v, ok := data.(T); ok {
		return v, nil
	}

	return *new(T), nil
}
