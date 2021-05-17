package ratelimit

import (
	"time"

	"github.com/juju/ratelimit"
)

type RateLimiter struct {
	limiter *ratelimit.Bucket
}

func NewRateLimiter(capacity int64) *RateLimiter {
	return &RateLimiter{
		limiter: ratelimit.NewBucket(time.Second/time.Duration(capacity), capacity),
	}
}

func (l *RateLimiter) Do(count int64) bool {
	return l.limiter.TakeAvailable(count) > 0
}

func (l *RateLimiter) Available() int64 {
	return l.limiter.Available()
}
