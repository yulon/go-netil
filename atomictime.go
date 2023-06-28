package netil

import (
	"sync/atomic"
	"time"
)

type AtomicTime int64

func (at *AtomicTime) Set(t time.Time) {
	atomic.StoreInt64((*int64)(at), t.Unix())
}

func (at *AtomicTime) Get() time.Time {
	return time.Unix(atomic.LoadInt64((*int64)(at)), 0)
}

func (at *AtomicTime) IsZero() bool {
	return atomic.LoadInt64((*int64)(at)) == 0
}

type AtomicDuration int64

func (ad *AtomicDuration) Set(d time.Duration) {
	atomic.StoreInt64((*int64)(ad), int64(d))
}

func (ad *AtomicDuration) Get() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(ad)))
}
