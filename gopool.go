package netil

type GoPool struct {
	ql *QuantityLimiter
	ch chan func()
	n  int
}

func NewGoPool(max int) *GoPool {
	return &GoPool{
		NewQuantityLimiter(uint64(max)),
		make(chan func()),
		0,
	}
}

func (gp *GoPool) Limit(max int) {
	gp.ql.LimitWith(uint64(max), func(n uint64) {
		c := gp.n - max
		if c <= 0 {
			return
		}
		for i := 0; i < c; i++ {
			gp.ch <- nil
			gp.n--
		}
	})
}

func (gp *GoPool) Go(f func()) {
	if gp == nil {
		go f()
		return
	}
	gp.ql.LockWith(func(n, max uint64) {
		if max == 0 {
			go func() {
				f()
				gp.ql.Unlock()
			}()
			return
		}
		if gp.n < int(n) {
			go func() {
				for {
					f := <-gp.ch
					if f == nil {
						return
					}
					f()
					gp.ql.Unlock()
				}
			}()
			gp.n++
		}
		gp.ch <- f
	})
}

var goPool = NewGoPool(0)

func LimitGoNum(max int) {
	goPool.Limit(max)
}

func Go(f func()) {
	goPool.Go(f)
}
