package ratelimit

import "time"

type Limiter struct {
	token          chan struct{}
	refillInterval time.Duration
	refillNumber   int
}

func MustNewRateLimiter(capacity int, refillInterval time.Duration, refillNumber int) *Limiter {
	l := &Limiter{
		token:          make(chan struct{}, capacity),
		refillInterval: refillInterval,
		refillNumber:   refillNumber,
	}
	// fill the token
	for i := 0; i < capacity; i++ {
		l.token <- struct{}{}
	}
	go l.refill()

	return l
}

func (l *Limiter) Acquire() error {
	select {
	case <-l.token:
		return nil
	default:
		return ErrLimitExceeded
	}
}

func (l *Limiter) refill() {
	for {
		time.Sleep(l.refillInterval)
		for i := 0; i < l.refillNumber; i++ {
			l.token <- struct{}{}
		}
	}
}
