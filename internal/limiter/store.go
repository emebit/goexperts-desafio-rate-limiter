package limiter

import (
	"context"
	"time"
)

// Interface para facilitar a utilização de outro mecanismo de persistência.
type Store interface {
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, duration time.Duration) error
}
