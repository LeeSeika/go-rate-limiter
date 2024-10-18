package ratelimit

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIPLimiter(t *testing.T) {
	Convey("TestIPLimiter", t, func() {

		Convey("Test Single IP Limit", func() {
			il := MustNewIPLimiter(3, 3, 1*time.Second)

			ip := generateRandomIPv4()

			exceedTimes := 0

			for i := 0; i < 25; i++ {
				err := il.TryAdd(ip)
				if err != nil {
					// t.Errorf("TryAdd failed: %v", err)
					So(i, ShouldEqual, 3+exceedTimes+3*exceedTimes)
					exceedTimes++
					time.Sleep(1 * time.Second)
				}
			}

		})
	})
}

func generateRandomIPv4() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}
