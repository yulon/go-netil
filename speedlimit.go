package netil

import (
	"sync"
	"sync/atomic"
	"time"
)

type SpeedValve struct {
	nPerS uint64
}

func NewSpeedValve(bytesPerSecond uint64) *SpeedValve {
	return &SpeedValve{nPerS: bytesPerSecond}
}

func (sv *SpeedValve) SetSpeed(bytesPerSecond uint64) {
	atomic.StoreUint64(&sv.nPerS, bytesPerSecond)
}

func (sv *SpeedValve) GetSpeed() uint64 {
	return atomic.LoadUint64(&sv.nPerS)
}

func (sv *SpeedValve) NewLimiter() *SpeedLimiter {
	return &SpeedLimiter{SpeedValve: sv}
}

type SpeedLimiter struct {
	*SpeedValve
	t   int64
	ct  int64
	n   uint64
	o   uint64
	mtx sync.Mutex
}

func NewSpeedLimiter(bytesPerSecond uint64) *SpeedLimiter {
	return &SpeedLimiter{SpeedValve: NewSpeedValve(bytesPerSecond)}
}

/*func (slr *SpeedLimiter) Wait(n int) {
	if slr == nil || slr.GetSpeed() == 0 {
		return
	}

	slr.mtx.Lock()
	defer slr.mtx.Unlock()

	nPerS := slr.GetSpeed()
	if nPerS == 0 {
		slr.n = 0
		slr.o = 0
		return
	}

	now := NowUnixMilli()

	n64 := uint64(n)

	ela := now - slr.t
	cela := now - slr.ct
	if ela > 1000 || cela > 1000 {
		slr.t = now
		slr.ct = now
		slr.n = n64
		slr.o = 0
		return
	}

	if slr.o > 0 {
		if slr.o >= n64 {
			slr.o -= n64
			return
		}
		n64 -= slr.o
		slr.o = 0
	}

	slr.n += n64

	if cela < 100 {
		return
	}

	slr.ct = now

	fPerS := float64(nPerS)
	expEla := int64(float64(slr.n) * 1000.0 / fPerS)
	delay := expEla - ela

	if delay > 30 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		slr.t = now
		slr.ct = now
		slr.n = 0
	} else if delay < -30 {
		owedDelay := 0 - delay
		slr.o = uint64(fPerS * float64(owedDelay) / 1000.0)
		slr.n = 0
	}
}*/

func (slr *SpeedLimiter) Wait(n int) {
	if slr == nil || slr.GetSpeed() == 0 {
		return
	}

	slr.mtx.Lock()
	defer slr.mtx.Unlock()

	nPerS := slr.GetSpeed()
	if nPerS == 0 {
		slr.n = 0
		slr.o = 0
		return
	}

	now := NowUnixMilli()

	n64 := uint64(n)
	if slr.o > 0 {
		if slr.o >= n64 {
			slr.o -= n64
			return
		}
		n64 -= slr.o
		slr.o = 0
	}

	if slr.n == 0 {
		slr.t = now
	}
	slr.n += n64

	if slr.n < kb4 {
		return
	}

	fPerS := float64(nPerS)
	ela := now - slr.t
	expEla := int64(float64(slr.n) * 1000.0 / fPerS)
	delay := expEla - ela

	if delay > 30 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		slr.n = 0
		return
	}

	if delay < 0 {
		owedDelay := 0 - delay
		slr.o = uint64(fPerS * (float64(owedDelay) / 1000.0))
		slr.n = 0
		return
	}

	if ela > 1000 {
		slr.n = 0
	}
}
