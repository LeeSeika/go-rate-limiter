package ratelimit

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	rl := MustNewRateLimiter(10, 1*time.Second, 5)

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			err := rl.Acquire()
			if err != nil {
				t.Errorf("Acquire failed: %v", err)
				t.Fail()
			}
		}
		time.Sleep(1 * time.Millisecond)
		for i := 0; i < 5; i++ {
			err := rl.Acquire()
			if err != nil {
				t.Errorf("Acquire failed: %v", err)
				t.Fail()
			}
		}
		wg.Done()
	}()

	wg.Wait()

}
