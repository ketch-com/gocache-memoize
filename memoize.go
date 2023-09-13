package memoize

import (
	"context"
	"github.com/eko/gocache/lib/v4/cache"
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
// cache, then the function is called and it's result is cached.
func (m *Memoizer[T]) Memoize(ctx context.Context, key string, fn func(context.Context) (T, error)) (value T, err error) {
	// Check cache
	value, err = m.Get(ctx, key)
	if err != nil {
		return
	}

	var data any
	data, err, _ = m.group.Do(key, func() (any, error) {
		value, err := fn(ctx)
		if err != nil {
			m.Set(ctx, key, value)
		}

		return value, err
	})
	if err != nil {
		return
	}

	if v, ok := data.(T); ok {
		value = v
	}

	return
}

