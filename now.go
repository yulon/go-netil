package netil

import (
	"sync/atomic"
	"time"
)

func NewNower(interval time.Duration, now func() int64) *int64 {
	n := new(int64)
	go func() {
		for {
			atomic.StoreInt64(n, now())
			time.Sleep(interval)
		}
	}()
	return n
}

var nowUnixMilliPtr = NewNower(50*time.Millisecond, func() int64 {
	return time.Now().UnixMilli()
})

func NowUnixMilli() int64 {
	return atomic.LoadInt64(nowUnixMilliPtr)
}

func Now() time.Time {
	return time.UnixMilli(NowUnixMilli())
}
