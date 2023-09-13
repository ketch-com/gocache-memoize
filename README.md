# gocache-memoize

Memoization implemented using [eko/gocache](github.com/eko/gocache).

## Installation

To begin working with the latest version of `gocache-memoize`, you can import the library in your project:

```
go get github.com/ketch-com/gocache-memoize
```

## Usage

```go
import (
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/ketch-com/gocache-memoize"
)

cacheManager := cache.New[[]byte](...)

m := memoize.NewMemoizer(cacheManager)
value, err := m.Memoize(ctx, "key1", func(ctx context.Context) (any, error) {
	return "test", nil
})
```
