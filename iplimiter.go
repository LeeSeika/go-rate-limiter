package ratelimit

import (
	"crypto/md5"
	"io"
	"sync"
	"time"
)

type IPLimiter struct {
	buckets          []map[string]int
	mutexs           []*sync.Mutex
	lastAcquireTimes []time.Time

	bucketCap   int
	limitNumber int
	windowSize  time.Duration
}

func MustNewIPLimiter(bucketCap int, limitNumber int, windowSize time.Duration) *IPLimiter {
	if bucketCap <= 0 || limitNumber <= 0 || windowSize <= 0 {
		panic("invalid arguments when creating IPLimiter")
	}
	buckets := make([]map[string]int, bucketCap)
	mutexs := make([]*sync.Mutex, bucketCap)
	lastAcquireTimes := make([]time.Time, bucketCap)
	for i := 0; i < bucketCap; i++ {
		buckets[i] = make(map[string]int)
		mutexs[i] = &sync.Mutex{}
		lastAcquireTimes[i] = time.Now()
	}

	return &IPLimiter{
		buckets:          buckets,
		mutexs:           mutexs,
		bucketCap:        bucketCap,
		limitNumber:      limitNumber,
		windowSize:       windowSize,
		lastAcquireTimes: lastAcquireTimes,
	}
}

func (m *IPLimiter) TryAdd(ip string) error {
	hash := md5.New()
	_, err := io.WriteString(hash, ip)
	if err != nil {
		return err
	}
	idx := int(hash.Sum(nil)[0]) % m.bucketCap

	m.mutexs[idx].Lock()
	defer m.mutexs[idx].Unlock()

	now := time.Now()
	bucket := m.buckets[idx]

	if now.Sub(m.lastAcquireTimes[idx]) > m.windowSize {
		m.lastAcquireTimes[idx] = now
		bucket[ip] = 0
	}

	reqCount := bucket[ip]
	if reqCount >= m.limitNumber {
		return ErrTooManyRequests
	}

	bucket[ip] = reqCount + 1
	return nil
}
