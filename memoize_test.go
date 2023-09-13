package memoize_test

import (
	"context"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	memoize "go.ketch.com/lib/gocache-memoize"
	"testing"
	"time"
)

type mockStore struct {
	Value any
	Error error
}

func (m *mockStore) Get(ctx context.Context, key any) (any, error) {
	return m.Value, m.Error
}

func (m *mockStore) GetWithTTL(ctx context.Context, key any) (any, time.Duration, error) {
	return m.Value, 10 * time.Second, m.Error
}

func (m *mockStore) Set(ctx context.Context, key any, value any, options ...store.Option) error {
	m.Value = value
	return nil
}

func (m *mockStore) Delete(ctx context.Context, key any) error {
	m.Value = nil
	return nil
}

func (m *mockStore) Invalidate(ctx context.Context, options ...store.InvalidateOption) error {
	return nil
}

func (m *mockStore) Clear(ctx context.Context) error {
	return nil
}

func (m *mockStore) GetType() string {
	return "mock"
}

func TestMemoize(t *testing.T) {
	c := &mockStore{}
	m := memoize.NewMemoizer[string](cache.New[string](c))

	c.Error = store.NotFound{}
	value, err := m.Memoize(context.Background(), "key1", func(ctx context.Context) (string, error) {
		return "value1", nil
	})
	if err != nil {
		t.Fatal(err)
	}

	c.Error = nil

	if value != "value1" {
		t.Fatal("value not matched")
	}

	value, err = m.Memoize(context.Background(), "key1", func(ctx context.Context) (string, error) {
		return "value2", nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if value != "value1" {
		t.Fatal("value not matched")
	}
}
