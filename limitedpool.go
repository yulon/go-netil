package netil

import (
	"container/list"
)

type LimitedPool struct {
	ql   *QuantityLimiter
	bufs *list.List
	new  func() any
}

func NewLimitedPool(max uint64, new func() any) *LimitedPool {
	return &LimitedPool{
		new:  new,
		ql:   NewQuantityLimiter(max),
		bufs: list.New(),
	}
}

func (lp *LimitedPool) Limit(max uint64) {
	lp.ql.LimitWith(max, func(n uint64) {
		if uint64(lp.bufs.Len()) <= max {
			return
		}
		delN := uint64(lp.bufs.Len()) - max
		it := lp.bufs.Front()
		for i := uint64(0); i < delN; i++ {
			next := it.Next()
			lp.bufs.Remove(it)
			it = next
		}
	})
}

func (lp *LimitedPool) Alloc() (val any) {
	lp.ql.LockWith(func(n, max uint64) {
		if lp.bufs.Len() == 0 {
			val = lp.new()
			return
		}
		it := lp.bufs.Front()
		val = it.Value
		lp.bufs.Remove(it)
	})
	return
}

func (lp *LimitedPool) Recycle(val any) {
	lp.ql.UnlockWith(func(n, max uint64) {
		if uint64(lp.bufs.Len()) >= max {
			return
		}
		lp.bufs.PushBack(val)
	})
}
