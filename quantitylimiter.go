package netil

import "sync"

type QuantityLimiter struct {
	max uint64
	c   uint64

	mtx  sync.Mutex
	cond *sync.Cond
}

func NewQuantityLimiter(max uint64) *QuantityLimiter {
	ql := &QuantityLimiter{max: max}
	ql.cond = sync.NewCond(&ql.mtx)
	return ql
}

func (ql *QuantityLimiter) Limit(max uint64) {
	ql.LimitWith(max, nil)
}

func (ql *QuantityLimiter) LimitWith(max uint64, f func(n uint64)) {
	ql.mtx.Lock()

	oldMax := ql.max
	ql.max = max

	if f != nil {
		f(ql.c)
	}

	if ql.c <= max && ql.c > oldMax && oldMax > 0 {
		ql.mtx.Unlock()
		ql.cond.Broadcast()
		return
	}

	ql.mtx.Unlock()
}

func (ql *QuantityLimiter) Lock() {
	ql.LockWith(nil)
}

func (ql *QuantityLimiter) LockWith(f func(n, max uint64)) {
	ql.mtx.Lock()
	defer ql.mtx.Unlock()

	ql.c++
	for ql.c > ql.max && ql.max > 0 {
		ql.cond.Wait()
	}

	if f != nil {
		f(ql.c, ql.max)
	}
}

func (ql *QuantityLimiter) Unlock() {
	ql.UnlockWith(nil)
}

func (ql *QuantityLimiter) UnlockWith(f func(n, max uint64)) {
	ql.mtx.Lock()

	if ql.c == 0 {
		panic(ql.c)
	}
	ql.c--

	if f != nil {
		f(ql.c, ql.max)
	}

	if ql.c == ql.max && ql.max > 0 {
		ql.mtx.Unlock()
		ql.cond.Broadcast()
	} else {
		ql.mtx.Unlock()
	}
}
