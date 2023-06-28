package netil

import (
	"net"
	"sync"
)

type PacketConnPool struct {
	pcons map[string]*PoolPacketConn
	max   int
	mtx   sync.Mutex
}

func NewPacketConnPool(max int) *PacketConnPool {
	return &PacketConnPool{max: max}
}

type PoolPacketConn struct {
	*EavesdroppedPacketConn

	vaddr       string
	listenToken string

	pool *PacketConnPool
}

func (pubCon *PoolPacketConn) closeLF() error {
	delete(pubCon.pool.pcons, pubCon.vaddr)
	return pubCon.EavesdroppedPacketConn.Close()
}

func (pubCon *PoolPacketConn) Close() error {
	return pubCon.pool.Del(pubCon.vaddr, pubCon.listenToken)
}

func (pool *PacketConnPool) Get(vaddr string, listenToken string, listenPubCon func() (net.PacketConn, error)) (pubCon *PoolPacketConn, err error) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	if pool.pcons == nil {
		pool.pcons = map[string]*PoolPacketConn{}
		pubCon = &PoolPacketConn{}
	} else {
		var has bool
		pubCon, has = pool.pcons[vaddr]
		if has {
			if pubCon.listenToken == listenToken {
				return
			}
			pubCon.EavesdroppedPacketConn.Close()
		} else {
			pubCon = &PoolPacketConn{}
		}
	}

	rawCon, err := ListenPacketOrDirect(listenPubCon)
	if err != nil {
		pubCon = nil
		return
	}
	pubCon.EavesdroppedPacketConn = EavesdropPacketConn(rawCon)

	if pool.max > 0 {
		for len(pool.pcons) >= pool.max {
			var olderOut *PoolPacketConn
			for _, o := range pool.pcons {
				if olderOut == nil || o.ReadWriteTime() < olderOut.ReadWriteTime() {
					olderOut = o
				}
			}
			olderOut.closeLF()
		}
	}

	pubCon.pool = pool
	pubCon.vaddr = vaddr
	pubCon.listenToken = listenToken
	pubCon.EnableReadWriteTime()
	pool.pcons[vaddr] = pubCon
	return
}

func (pool *PacketConnPool) delLF(vaddr string, listenToken string) error {
	if pool.pcons == nil {
		return nil
	}

	pubCon, has := pool.pcons[vaddr]
	if !has || pubCon.listenToken != listenToken {
		return nil
	}

	return pubCon.closeLF()
}

func (pool *PacketConnPool) Del(vaddr string, listenToken string) error {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	return pool.delLF(vaddr, listenToken)
}

func (pool *PacketConnPool) DelByListenToken(listenToken string) {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	if pool.pcons == nil {
		return
	}

	for _, pubCon := range pool.pcons {
		if pubCon.listenToken == listenToken {
			pubCon.closeLF()
		}
	}
}

func (pool *PacketConnPool) Close() error {
	pool.mtx.Lock()
	defer pool.mtx.Unlock()

	if pool.pcons == nil {
		return nil
	}

	for _, pubCon := range pool.pcons {
		pubCon.EavesdroppedPacketConn.Close()
	}
	pool.pcons = nil
	return nil
}
