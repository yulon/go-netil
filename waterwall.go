package netil

import (
	"net"
	"sync"
)

type waterWallIPInf struct {
	ErrCnt       int
	BlockEndTime int64
}

type WaterWall struct {
	ipInfs      map[string]*waterWallIPInf
	userErrCnts map[string]int
	mtx         sync.RWMutex
}

func NewWaterWall() *WaterWall {
	return &WaterWall{
		ipInfs:      map[string]*waterWallIPInf{},
		userErrCnts: map[string]int{},
	}
}

func (ww *WaterWall) IsBlock(ip net.IP) bool {
	if ww == nil || len(ip) == 0 {
		return false
	}

	ipStr := NarrowIP(ip).String()

	ww.mtx.RLock()
	defer ww.mtx.RUnlock()

	ipInf, has := ww.ipInfs[ipStr]
	return has && ipInf.BlockEndTime != 0 && ipInf.BlockEndTime > NowUnixMilli()
}

func (ww *WaterWall) IsBlockAddr(addr net.Addr) bool {
	return ww.IsBlock(IPFromAddr(addr))
}

func (ww *WaterWall) Fail(ip net.IP, user string) bool {
	if ww == nil || len(ip) == 0 {
		return true
	}

	ipStr := NarrowIP(ip).String()

	ww.mtx.Lock()
	defer ww.mtx.Unlock()

	ipInf, has := ww.ipInfs[ipStr]
	if !has {
		ipInf = &waterWallIPInf{ErrCnt: 1}
		ww.ipInfs[ipStr] = ipInf
	} else if ipInf.ErrCnt < 6 {
		ipInf.ErrCnt++
	}
	ec := ipInf.ErrCnt

	if len(user) != 0 {
		uec := ww.userErrCnts[user]
		if uec < 6 {
			uec++
			ww.userErrCnts[user] = uec
		}
		if uec > ec {
			ec = uec
		}
	}

	switch ec {
	case 1:
		return true
	case 2:
		return false
	case 3:
		ipInf.BlockEndTime = NowUnixMilli() + 60*1000
	case 4:
		ipInf.BlockEndTime = NowUnixMilli() + 10*60*1000
	case 5:
		ipInf.BlockEndTime = NowUnixMilli() + 30*60*1000
	case 6:
		ipInf.BlockEndTime = NowUnixMilli() + 60*60*1000
	default:
		panic(ec)
	}
	return false
}

func (ww *WaterWall) FailAddr(addr net.Addr, user string) bool {
	return ww.Fail(IPFromAddr(addr), user)
}

func (ww *WaterWall) Success(ip net.IP, user string) {
	if ww == nil || len(ip) == 0 {
		return
	}

	ipStr := NarrowIP(ip).String()

	ww.mtx.Lock()
	defer ww.mtx.Unlock()

	delete(ww.ipInfs, ipStr)
	delete(ww.userErrCnts, user)
}

func (ww *WaterWall) SuccessAddr(addr net.Addr, user string) {
	ww.Success(IPFromAddr(addr), user)
}
