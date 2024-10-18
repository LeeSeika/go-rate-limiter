package ratelimit

import "errors"

var (
	ErrLimitExceeded   = errors.New("rate limit exceeded")
	ErrTooManyRequests = errors.New("too many requests")
)
