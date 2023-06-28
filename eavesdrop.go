package netil

import (
	"net"
	"sync"
	"sync/atomic"
)

type Eavesdropper struct {
	onAnyReadBefore  func(bn int)
	onAnyReadAfter   func(data []byte)
	onAnyWriteBefore func(data []byte)
	onAnyWriteAfter  func(wn int)

	onStreamReadBefore  func(bn int)
	onStreamReadAfter   func(data []byte)
	onStreamWriteBefore func(data []byte)
	onStreamWriteAfter  func(wn int)

	onPacketReadBefore  func(bn int)
	onPacketReadAfter   func(data []byte, from net.Addr)
	onPacketWriteBefore func(data []byte, to net.Addr)
	onPacketWriteAfter  func(wn int)

	onClose   func()
	closeOnce sync.Once

	rTime int64
	wTime int64
}

func (edr *Eavesdropper) OnAnyWriteBefore(onAnyWriteBefore func(data []byte)) {
	if edr.onAnyWriteBefore == nil {
		edr.onAnyWriteBefore = onAnyWriteBefore
		return
	}
	oldOnAnyWriteBefore := edr.onAnyWriteBefore
	edr.onAnyWriteBefore = func(data []byte) {
		oldOnAnyWriteBefore(data)
		onAnyWriteBefore(data)
	}
}

func (edr *Eavesdropper) NotifyAnyWriteBefore(data []byte) {
	if edr == nil {
		return
	}
	if edr.onAnyWriteBefore != nil {
		edr.onAnyWriteBefore(data)
	}
}

func (edr *Eavesdropper) OnAnyWriteAfter(onAnyWriteAfter func(wn int)) {
	if edr.onAnyWriteAfter == nil {
		edr.onAnyWriteAfter = onAnyWriteAfter
		return
	}
	oldOnAnyWriteAfter := edr.onAnyWriteAfter
	edr.onAnyWriteAfter = func(wn int) {
		oldOnAnyWriteAfter(wn)
		onAnyWriteAfter(wn)
	}
}

func (edr *Eavesdropper) NotifyAnyWriteAfter(wn int) {
	if edr == nil {
		return
	}
	if edr.onAnyWriteAfter != nil {
		edr.onAnyWriteAfter(wn)
	}
}

func (edr *Eavesdropper) OnAnyReadBefore(onAnyReadBefore func(bn int)) {
	if edr.onAnyReadBefore == nil {
		edr.onAnyReadBefore = onAnyReadBefore
		return
	}
	oldOnAnyReadBefore := edr.onAnyReadBefore
	edr.onAnyReadBefore = func(bn int) {
		oldOnAnyReadBefore(bn)
		onAnyReadBefore(bn)
	}
}

func (edr *Eavesdropper) NotifyAnyReadBefore(bn int) {
	if edr == nil {
		return
	}
	if edr.onAnyReadBefore != nil {
		edr.onAnyReadBefore(bn)
	}
}

func (edr *Eavesdropper) OnAnyReadAfter(onAnyReadAfter func(data []byte)) {
	if edr.onAnyReadAfter == nil {
		edr.onAnyReadAfter = onAnyReadAfter
		return
	}
	oldOnAnyReadAfter := edr.onAnyReadAfter
	edr.onAnyReadAfter = func(data []byte) {
		oldOnAnyReadAfter(data)
		onAnyReadAfter(data)
	}
}

func (edr *Eavesdropper) NotifyAnyReadAfter(data []byte) {
	if edr == nil {
		return
	}
	if edr.onAnyReadAfter != nil {
		edr.onAnyReadAfter(data)
	}
}

func (edr *Eavesdropper) OnStreamWriteBefore(onStreamWriteBefore func(data []byte)) {
	if edr.onStreamWriteBefore == nil {
		edr.onStreamWriteBefore = onStreamWriteBefore
		return
	}
	oldOnStreamWriteBefore := edr.onStreamWriteBefore
	edr.onStreamWriteBefore = func(data []byte) {
		oldOnStreamWriteBefore(data)
		onStreamWriteBefore(data)
	}
}

func (edr *Eavesdropper) NotifyStreamWriteBefore(data []byte) {
	if edr == nil {
		return
	}
	if edr.onAnyWriteBefore != nil {
		edr.onAnyWriteBefore(data)
	}
	if edr.onStreamWriteBefore != nil {
		edr.onStreamWriteBefore(data)
	}
}

func (edr *Eavesdropper) OnStreamWriteAfter(onStreamWriteAfter func(wn int)) {
	if edr.onStreamWriteAfter == nil {
		edr.onStreamWriteAfter = onStreamWriteAfter
		return
	}
	oldOnStreamWriteAfter := edr.onStreamWriteAfter
	edr.onStreamWriteAfter = func(wn int) {
		oldOnStreamWriteAfter(wn)
		onStreamWriteAfter(wn)
	}
}

func (edr *Eavesdropper) NotifyStreamWriteAfter(wn int) {
	if edr == nil {
		return
	}
	if edr.onAnyWriteAfter != nil {
		edr.onAnyWriteAfter(wn)
	}
	if edr.onStreamWriteAfter != nil {
		edr.onStreamWriteAfter(wn)
	}
}

func (edr *Eavesdropper) OnStreamReadBefore(onStreamReadBefore func(bn int)) {
	if edr.onStreamReadBefore == nil {
		edr.onStreamReadBefore = onStreamReadBefore
		return
	}
	oldOnStreamReadBefore := edr.onStreamReadBefore
	edr.onStreamReadBefore = func(bn int) {
		oldOnStreamReadBefore(bn)
		onStreamReadBefore(bn)
	}
}

func (edr *Eavesdropper) NotifyStreamReadBefore(bn int) {
	if edr == nil {
		return
	}
	if edr.onStreamReadBefore != nil {
		edr.onStreamReadBefore(bn)
	}
	if edr.onAnyReadBefore != nil {
		edr.onAnyReadBefore(bn)
	}
}

func (edr *Eavesdropper) OnStreamReadAfter(onStreamReadAfter func(data []byte)) {
	if edr.onStreamReadAfter == nil {
		edr.onStreamReadAfter = onStreamReadAfter
		return
	}
	oldOnStreamReadAfter := edr.onStreamReadAfter
	edr.onStreamReadAfter = func(data []byte) {
		oldOnStreamReadAfter(data)
		onStreamReadAfter(data)
	}
}

func (edr *Eavesdropper) NotifyStreamReadAfter(data []byte) {
	if edr == nil {
		return
	}
	if edr.onStreamReadAfter != nil {
		edr.onStreamReadAfter(data)
	}
	if edr.onAnyReadAfter != nil {
		edr.onAnyReadAfter(data)
	}
}

func (edr *Eavesdropper) OnPacketWriteBefore(onPacketWriteBefore func(data []byte, to net.Addr)) {
	if edr.onPacketWriteBefore == nil {
		edr.onPacketWriteBefore = onPacketWriteBefore
		return
	}
	oldOnPacketWriteBefore := edr.onPacketWriteBefore
	edr.onPacketWriteBefore = func(data []byte, to net.Addr) {
		oldOnPacketWriteBefore(data, to)
		onPacketWriteBefore(data, to)
	}
}

func (edr *Eavesdropper) NotifyPacketWriteBefore(data []byte, to net.Addr) {
	if edr == nil {
		return
	}
	if edr.onAnyWriteBefore != nil {
		edr.onAnyWriteBefore(data)
	}
	if edr.onPacketWriteBefore != nil {
		edr.onPacketWriteBefore(data, to)
	}
}

func (edr *Eavesdropper) OnPacketWriteAfter(onPacketWriteAfter func(wn int)) {
	if edr.onPacketWriteAfter == nil {
		edr.onPacketWriteAfter = onPacketWriteAfter
		return
	}
	oldOnPacketWriteAfter := edr.onPacketWriteAfter
	edr.onPacketWriteAfter = func(wn int) {
		oldOnPacketWriteAfter(wn)
		onPacketWriteAfter(wn)
	}
}

func (edr *Eavesdropper) NotifyPacketWriteAfter(wn int) {
	if edr == nil {
		return
	}
	if edr.onAnyWriteAfter != nil {
		edr.onAnyWriteAfter(wn)
	}
	if edr.onPacketWriteAfter != nil {
		edr.onPacketWriteAfter(wn)
	}
}

func (edr *Eavesdropper) OnPacketReadBefore(onPacketReadBefore func(bn int)) {
	if edr.onPacketReadBefore == nil {
		edr.onPacketReadBefore = onPacketReadBefore
		return
	}
	oldOnPacketReadBefore := edr.onPacketReadBefore
	edr.onPacketReadBefore = func(bn int) {
		oldOnPacketReadBefore(bn)
		onPacketReadBefore(bn)
	}
}

func (edr *Eavesdropper) NotifyPacketReadBefore(bn int) {
	if edr == nil {
		return
	}
	if edr.onPacketReadBefore != nil {
		edr.onPacketReadBefore(bn)
	}
	if edr.onAnyReadBefore != nil {
		edr.onAnyReadBefore(bn)
	}
}

func (edr *Eavesdropper) OnPacketReadAfter(onPacketReadAfter func(data []byte, from net.Addr)) {
	if edr.onPacketReadAfter == nil {
		edr.onPacketReadAfter = onPacketReadAfter
		return
	}
	oldOnPacketReadAfter := edr.onPacketReadAfter
	edr.onPacketReadAfter = func(data []byte, from net.Addr) {
		oldOnPacketReadAfter(data, from)
		onPacketReadAfter(data, from)
	}
}

func (edr *Eavesdropper) NotifyPacketReadAfter(data []byte, from net.Addr) {
	if edr == nil {
		return
	}
	if edr.onPacketReadAfter != nil {
		edr.onPacketReadAfter(data, from)
	}
	if edr.onAnyReadAfter != nil {
		edr.onAnyReadAfter(data)
	}
}

func (edr *Eavesdropper) OnClose(onClose func()) {
	if edr.onClose == nil {
		edr.onClose = onClose
		return
	}
	oldOnClose := edr.onClose
	edr.onClose = func() {
		oldOnClose()
		onClose()
	}
}

func (edr *Eavesdropper) NotifyClose() {
	if edr == nil || edr.onClose == nil {
		return
	}
	edr.closeOnce.Do(func() {
		edr.onClose()
	})
}

func (edr *Eavesdropper) LimitReadSpeed(limiter *SpeedLimiter) bool {
	if limiter == nil {
		return false
	}
	edr.OnAnyReadAfter(func(data []byte) {
		limiter.Wait(len(data))
	})
	/*bufSz := new(int)
	edr.OnAnyReadBefore(func(bn int) {
		*bufSz = bn
		limiter.Wait(0)
	})
	edr.OnAnyReadAfter(func(data []byte) {
		limiter.Fixed(len(data) - *bufSz)
	})*/
	return true
}

func (edr *Eavesdropper) LimitWriteSpeed(limiter *SpeedLimiter) bool {
	if limiter == nil {
		return false
	}
	edr.OnAnyWriteBefore(func(data []byte) {
		limiter.Wait(len(data))
	})
	/*dataSz := new(int)
	edr.OnAnyWriteBefore(func(data []byte) {
		*dataSz = len(data)
		limiter.Wait(*dataSz)
	})
	edr.OnAnyWriteAfter(func(wn int) {
		limiter.Fixed(wn - *dataSz)
	})*/
	return true
}

func (edr *Eavesdropper) EnableReadWriteTime() bool {
	if atomic.LoadInt64(&edr.rTime) != 0 && atomic.LoadInt64(&edr.wTime) != 0 {
		return false
	}

	edr.UpdateReadWriteTime()

	edr.OnAnyReadAfter(func([]byte) {
		edr.UpdateReadTime()
	})
	edr.OnAnyWriteBefore(func([]byte) {
		edr.UpdateWriteTime()
	})
	edr.OnClose(func() {
		edr.UpdateReadWriteTime()
	})
	return true
}

func (edr *Eavesdropper) UpdateReadTime() {
	atomic.StoreInt64(&edr.rTime, NowUnixMilli())
}

func (edr *Eavesdropper) UpdateWriteTime() {
	atomic.StoreInt64(&edr.rTime, NowUnixMilli())
}

func (edr *Eavesdropper) UpdateReadWriteTime() {
	t := NowUnixMilli()
	atomic.StoreInt64(&edr.rTime, t)
	atomic.StoreInt64(&edr.wTime, t)
}

func (edr *Eavesdropper) ReadTime() int64 {
	return atomic.LoadInt64(&edr.rTime)
}

func (edr *Eavesdropper) WriteTime() int64 {
	return atomic.LoadInt64(&edr.wTime)
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func (edr *Eavesdropper) ReadWriteTime() int64 {
	return max(edr.ReadTime(), edr.WriteTime())
}

type EavesdroppedConn struct {
	*Eavesdropper
	net.Conn
}

func EavesdropConn(con net.Conn) *EavesdroppedConn {
	if con == nil {
		return nil
	}
	ec, isEC := con.(*EavesdroppedConn)
	if isEC {
		return ec
	}
	return &EavesdroppedConn{&Eavesdropper{}, con}
}

func (edr *Eavesdropper) EavesdropConn(con net.Conn) *EavesdroppedConn {
	if con == nil {
		return nil
	}

	ec, isEC := con.(*EavesdroppedConn)
	if !isEC {
		return &EavesdroppedConn{edr, con}
	}

	ecEdr := ec.GetEavesdropper()
	if edr.onAnyReadBefore != nil {
		ecEdr.OnAnyReadBefore(edr.onAnyReadBefore)
	}
	if edr.onAnyReadAfter != nil {
		ecEdr.OnAnyReadAfter(edr.onAnyReadAfter)
	}
	if edr.onAnyWriteBefore != nil {
		ecEdr.OnAnyWriteBefore(edr.onAnyWriteBefore)
	}
	if edr.onAnyWriteAfter != nil {
		ecEdr.OnAnyWriteAfter(edr.onAnyWriteAfter)
	}
	if edr.onStreamReadBefore != nil {
		ecEdr.OnStreamReadBefore(edr.onStreamReadBefore)
	}
	if edr.onStreamReadAfter != nil {
		ecEdr.OnStreamReadAfter(edr.onStreamReadAfter)
	}
	if edr.onStreamWriteBefore != nil {
		ecEdr.OnStreamWriteBefore(edr.onStreamWriteBefore)
	}
	if edr.onStreamWriteAfter != nil {
		ecEdr.OnStreamWriteAfter(edr.onStreamWriteAfter)
	}
	if edr.onPacketReadBefore != nil {
		ecEdr.OnPacketReadBefore(edr.onPacketReadBefore)
	}
	if edr.onPacketReadAfter != nil {
		ecEdr.OnPacketReadAfter(edr.onPacketReadAfter)
	}
	if edr.onPacketWriteBefore != nil {
		ecEdr.OnPacketWriteBefore(edr.onPacketWriteBefore)
	}
	if edr.onPacketWriteAfter != nil {
		ecEdr.OnPacketWriteAfter(edr.onPacketWriteAfter)
	}
	if edr.onClose != nil {
		ecEdr.OnClose(edr.onClose)
	}
	return ec
}

func (edr *Eavesdropper) TryEavesdropConn(con net.Conn) net.Conn {
	if con == nil || edr == nil {
		return con
	}

	ec, isEC := con.(*EavesdroppedConn)
	if !isEC {
		if edr.onAnyReadBefore == nil &&
			edr.onAnyReadAfter == nil &&
			edr.onAnyWriteBefore == nil &&
			edr.onAnyWriteAfter == nil &&
			edr.onStreamReadBefore == nil &&
			edr.onStreamReadAfter == nil &&
			edr.onStreamWriteBefore == nil &&
			edr.onStreamWriteAfter == nil &&
			edr.onPacketReadBefore == nil &&
			edr.onPacketReadAfter == nil &&
			edr.onPacketWriteBefore == nil &&
			edr.onPacketWriteAfter == nil &&
			edr.onClose == nil {
			return con
		}
		return &EavesdroppedConn{edr, con}
	}

	ecEdr := ec.GetEavesdropper()
	if edr.onAnyReadBefore != nil {
		ecEdr.OnAnyReadBefore(edr.onAnyReadBefore)
	}
	if edr.onAnyReadAfter != nil {
		ecEdr.OnAnyReadAfter(edr.onAnyReadAfter)
	}
	if edr.onAnyWriteBefore != nil {
		ecEdr.OnAnyWriteBefore(edr.onAnyWriteBefore)
	}
	if edr.onAnyWriteAfter != nil {
		ecEdr.OnAnyWriteAfter(edr.onAnyWriteAfter)
	}
	if edr.onStreamReadBefore != nil {
		ecEdr.OnStreamReadBefore(edr.onStreamReadBefore)
	}
	if edr.onStreamReadAfter != nil {
		ecEdr.OnStreamReadAfter(edr.onStreamReadAfter)
	}
	if edr.onStreamWriteBefore != nil {
		ecEdr.OnStreamWriteBefore(edr.onStreamWriteBefore)
	}
	if edr.onStreamWriteAfter != nil {
		ecEdr.OnStreamWriteAfter(edr.onStreamWriteAfter)
	}
	if edr.onPacketReadBefore != nil {
		ecEdr.OnPacketReadBefore(edr.onPacketReadBefore)
	}
	if edr.onPacketReadAfter != nil {
		ecEdr.OnPacketReadAfter(edr.onPacketReadAfter)
	}
	if edr.onPacketWriteBefore != nil {
		ecEdr.OnPacketWriteBefore(edr.onPacketWriteBefore)
	}
	if edr.onPacketWriteAfter != nil {
		ecEdr.OnPacketWriteAfter(edr.onPacketWriteAfter)
	}
	if edr.onClose != nil {
		ecEdr.OnClose(edr.onClose)
	}
	return ec
}

func (ec *EavesdroppedConn) GetEavesdropper() *Eavesdropper {
	return ec.Eavesdropper
}

func (ec *EavesdroppedConn) Read(buf []byte) (int, error) {
	ec.NotifyStreamReadBefore(len(buf))
	rn, err := ec.Conn.Read(buf)
	if rn > 0 {
		ec.NotifyStreamReadAfter(buf[:rn])
	}
	return rn, err
}

func (ec *EavesdroppedConn) Write(data []byte) (int, error) {
	ec.NotifyStreamWriteBefore(data)
	wn, err := ec.Conn.Write(data)
	if wn > 0 {
		ec.NotifyStreamWriteAfter(wn)
	}
	return wn, err
}

func (ec *EavesdroppedConn) Close() error {
	err := ec.Conn.Close()
	ec.NotifyClose()
	return err
}

type EavesdroppedPacketConn struct {
	*Eavesdropper
	net.PacketConn
}

func EavesdropPacketConn(pcon net.PacketConn) *EavesdroppedPacketConn {
	if pcon == nil {
		return nil
	}
	epc, isEPC := pcon.(*EavesdroppedPacketConn)
	if isEPC {
		return epc
	}
	return &EavesdroppedPacketConn{&Eavesdropper{}, pcon}
}

func (edr *Eavesdropper) EavesdropPacketConn(pcon net.PacketConn) *EavesdroppedPacketConn {
	if pcon == nil {
		return nil
	}

	epc, isEPC := pcon.(*EavesdroppedPacketConn)
	if !isEPC {
		return &EavesdroppedPacketConn{edr, pcon}
	}

	ecEdr := epc.GetEavesdropper()
	if edr.onAnyReadBefore != nil {
		ecEdr.OnAnyReadBefore(edr.onAnyReadBefore)
	}
	if edr.onAnyReadAfter != nil {
		ecEdr.OnAnyReadAfter(edr.onAnyReadAfter)
	}
	if edr.onAnyWriteBefore != nil {
		ecEdr.OnAnyWriteBefore(edr.onAnyWriteBefore)
	}
	if edr.onAnyWriteAfter != nil {
		ecEdr.OnAnyWriteAfter(edr.onAnyWriteAfter)
	}
	if edr.onStreamReadBefore != nil {
		ecEdr.OnStreamReadBefore(edr.onStreamReadBefore)
	}
	if edr.onStreamReadAfter != nil {
		ecEdr.OnStreamReadAfter(edr.onStreamReadAfter)
	}
	if edr.onStreamWriteBefore != nil {
		ecEdr.OnStreamWriteBefore(edr.onStreamWriteBefore)
	}
	if edr.onStreamWriteAfter != nil {
		ecEdr.OnStreamWriteAfter(edr.onStreamWriteAfter)
	}
	if edr.onPacketReadBefore != nil {
		ecEdr.OnPacketReadBefore(edr.onPacketReadBefore)
	}
	if edr.onPacketReadAfter != nil {
		ecEdr.OnPacketReadAfter(edr.onPacketReadAfter)
	}
	if edr.onPacketWriteBefore != nil {
		ecEdr.OnPacketWriteBefore(edr.onPacketWriteBefore)
	}
	if edr.onPacketWriteAfter != nil {
		ecEdr.OnPacketWriteAfter(edr.onPacketWriteAfter)
	}
	if edr.onClose != nil {
		ecEdr.OnClose(edr.onClose)
	}
	return epc
}

func (edr *Eavesdropper) TryEavesdropPacketConn(pcon net.PacketConn) net.PacketConn {
	if pcon == nil || edr == nil {
		return pcon
	}

	epc, isEPC := pcon.(*EavesdroppedPacketConn)
	if !isEPC {
		if edr.onAnyReadBefore == nil &&
			edr.onAnyReadAfter == nil &&
			edr.onAnyWriteBefore == nil &&
			edr.onAnyWriteAfter == nil &&
			edr.onStreamReadBefore == nil &&
			edr.onStreamReadAfter == nil &&
			edr.onStreamWriteBefore == nil &&
			edr.onStreamWriteAfter == nil &&
			edr.onPacketReadBefore == nil &&
			edr.onPacketReadAfter == nil &&
			edr.onPacketWriteBefore == nil &&
			edr.onPacketWriteAfter == nil &&
			edr.onClose == nil {
			return pcon
		}
		return &EavesdroppedPacketConn{edr, pcon}
	}

	ecEdr := epc.GetEavesdropper()
	if edr.onAnyReadBefore != nil {
		ecEdr.OnAnyReadBefore(edr.onAnyReadBefore)
	}
	if edr.onAnyReadAfter != nil {
		ecEdr.OnAnyReadAfter(edr.onAnyReadAfter)
	}
	if edr.onAnyWriteBefore != nil {
		ecEdr.OnAnyWriteBefore(edr.onAnyWriteBefore)
	}
	if edr.onAnyWriteAfter != nil {
		ecEdr.OnAnyWriteAfter(edr.onAnyWriteAfter)
	}
	if edr.onStreamReadBefore != nil {
		ecEdr.OnStreamReadBefore(edr.onStreamReadBefore)
	}
	if edr.onStreamReadAfter != nil {
		ecEdr.OnStreamReadAfter(edr.onStreamReadAfter)
	}
	if edr.onStreamWriteBefore != nil {
		ecEdr.OnStreamWriteBefore(edr.onStreamWriteBefore)
	}
	if edr.onStreamWriteAfter != nil {
		ecEdr.OnStreamWriteAfter(edr.onStreamWriteAfter)
	}
	if edr.onPacketReadBefore != nil {
		ecEdr.OnPacketReadBefore(edr.onPacketReadBefore)
	}
	if edr.onPacketReadAfter != nil {
		ecEdr.OnPacketReadAfter(edr.onPacketReadAfter)
	}
	if edr.onPacketWriteBefore != nil {
		ecEdr.OnPacketWriteBefore(edr.onPacketWriteBefore)
	}
	if edr.onPacketWriteAfter != nil {
		ecEdr.OnPacketWriteAfter(edr.onPacketWriteAfter)
	}
	if edr.onClose != nil {
		ecEdr.OnClose(edr.onClose)
	}
	return epc
}

func (epc *EavesdroppedPacketConn) GetEavesdropper() *Eavesdropper {
	return epc.Eavesdropper
}

func (epc *EavesdroppedPacketConn) ReadFrom(buf []byte) (int, net.Addr, error) {
	epc.NotifyPacketReadBefore(len(buf))
	rn, addr, err := epc.PacketConn.ReadFrom(buf)
	if rn > 0 {
		epc.NotifyPacketReadAfter(buf[:rn], addr)
	}
	return rn, addr, err
}

func (epc *EavesdroppedPacketConn) WriteTo(data []byte, addr net.Addr) (int, error) {
	epc.NotifyPacketWriteBefore(data, addr)
	wn, err := epc.PacketConn.WriteTo(data, addr)
	if wn > 0 {
		epc.NotifyPacketWriteAfter(wn)
	}
	return wn, err
}

func (epc *EavesdroppedPacketConn) Close() error {
	err := epc.PacketConn.Close()
	epc.NotifyClose()
	return err
}

/*type timeoutCloser struct {
	typ    byte
	expDur int64
	expTm  int64
	once   bool
	edr    *Eavesdropper
	closer io.Closer
}

var timeoutCloserCh chan *timeoutCloser
var timeoutCloserOnce sync.Once

var errCloseIdle = errors.New("force close idle connection")

func startTimeoutClosers() {
	timeoutCloserOnce.Do(func() {
		timeoutCloserCh = make(chan *timeoutCloser, 4096)
		go func() {
			tcs := list.New()
			waitTi := NowUnixMilli()
			waitDur := int64(30000)
			ticker := time.Tick(time.Duration(waitDur) * time.Millisecond)
			for {
				now := NowUnixMilli()
				select {
				case addTC := <-timeoutCloserCh:
					added := false
					for it := tcs.Front(); it != nil; it = it.Next() {
						tc := it.Value.(*timeoutCloser)
						if tc.expTm > addTC.expTm {
							tcs.InsertAfter(addTC, it)
							break
						}
					}
					if !added {
						tcs.PushBack(addTC)
					}
					waitEnd := waitTi + waitDur
					if addTC.expTm < waitEnd {
						waitEnd = addTC.expTm
					}
					waitDur = waitEnd - now
				case <-ticker:
					waitDur = int64(30000)
					for it := tcs.Front(); it != nil; {
						tc := it.Value.(*timeoutCloser)

						if now < tc.expTm {
							waitDur = tc.expTm - now
							break
						}

						var t int64
						switch tc.typ {
						case 0:
							t = tc.edr.ReadWriteTime()
						case 'r':
							t = tc.edr.ReadTime()
						case 'w':
							t = tc.edr.WriteTime()
						}

						next := it.Next()
						tc.expTm = t + tc.expDur
						if now >= tc.expTm {
							tc.closer.Close()
							tcs.Remove(it)
						} else if tc.once {
							tcs.Remove(it)
						} else {
							added := false
							for jt := next; jt != nil; jt = jt.Next() {
								if jt.Value.(*timeoutCloser).expTm > tc.expTm {
									tcs.MoveAfter(it, jt)
									if jt == next {
										next = jt.Prev()
									}
									added = true
									break
								}
							}
							if !added {
								tcs.MoveToBack(it)
							}
						}
						it = next
					}
				}
				if waitDur < 1000 {
					waitDur = 1000
				}
				ticker = time.Tick(time.Duration(waitDur) * time.Millisecond)
				waitTi = now
			}
		}()
	})
}

func (ec *EavesdroppedConn) SetIdleTimeout(expired time.Duration, once bool) {
	if ec == nil || ec.Conn == nil || expired == 0 {
		return
	}

	just := ec.EnableReadWriteTime()

	t := ec.ReadWriteTime()
	expDur := int64(expired / time.Millisecond)
	expTm := t + expDur
	now := NowUnixMilli()

	if !just && now >= expTm {
		ec.Close()
		return
	}

	if once && ec.Conn.SetDeadline(time.UnixMilli(now).Add(5*time.Second)) == nil {
		return
	}

	startTimeoutClosers()
	timeoutCloserCh <- &timeoutCloser{
		typ:    0,
		expDur: expDur,
		expTm:  expTm,
		once:   once,
		edr:    ec.Eavesdropper,
		closer: ec}
}

func (ec *EavesdroppedConn) SetReadTimeout(expired time.Duration, once bool) {
	if ec == nil || ec.Conn == nil || expired == 0 {
		return
	}

	just := ec.EnableReadWriteTime()

	t := ec.ReadTime()
	expDur := int64(expired / time.Millisecond)
	expTm := t + expDur
	now := NowUnixMilli()

	if !just && now >= expTm {
		ec.Close()
		return
	}

	if once && ec.Conn.SetReadDeadline(time.UnixMilli(now).Add(5*time.Second)) == nil {
		return
	}

	startTimeoutClosers()
	timeoutCloserCh <- &timeoutCloser{
		typ:    'r',
		expDur: expDur,
		expTm:  expTm,
		once:   once,
		edr:    ec.Eavesdropper,
		closer: ec}
}

func (ec *EavesdroppedConn) SetWriteTimeout(expired time.Duration, once bool) {
	if ec == nil || ec.Conn == nil || expired == 0 {
		return
	}

	just := ec.EnableReadWriteTime()

	t := ec.WriteTime()
	expDur := int64(expired / time.Millisecond)
	expTm := t + expDur
	now := NowUnixMilli()

	if !just && now >= expTm {
		ec.Close()
		return
	}

	if once && ec.Conn.SetWriteDeadline(time.UnixMilli(now).Add(5*time.Second)) == nil {
		return
	}

	startTimeoutClosers()
	timeoutCloserCh <- &timeoutCloser{
		typ:    'w',
		expDur: expDur,
		expTm:  expTm,
		once:   once,
		edr:    ec.Eavesdropper,
		closer: ec}
}

func (epc *EavesdroppedPacketConn) SetIdleTimeout(expired time.Duration, once bool) {
	if epc == nil || epc.PacketConn == nil || expired == 0 {
		return
	}

	just := epc.EnableReadWriteTime()

	t := epc.ReadWriteTime()
	expDur := int64(expired / time.Millisecond)
	expTm := t + expDur
	now := NowUnixMilli()

	if !just && now >= expTm {
		epc.Close()
		return
	}

	if once && epc.PacketConn.SetDeadline(time.UnixMilli(now).Add(5*time.Second)) == nil {
		return
	}

	startTimeoutClosers()
	timeoutCloserCh <- &timeoutCloser{
		typ:    0,
		expDur: expDur,
		expTm:  expTm,
		once:   once,
		edr:    epc.Eavesdropper,
		closer: epc}
}

func (epc *EavesdroppedPacketConn) SetReadTimeout(expired time.Duration, once bool) {
	if epc == nil || epc.PacketConn == nil || expired == 0 {
		return
	}

	just := epc.EnableReadWriteTime()

	t := epc.ReadTime()
	expDur := int64(expired / time.Millisecond)
	expTm := t + expDur
	now := NowUnixMilli()

	if !just && now >= expTm {
		epc.Close()
		return
	}

	if once && epc.PacketConn.SetReadDeadline(time.UnixMilli(now).Add(5*time.Second)) == nil {
		return
	}

	startTimeoutClosers()
	timeoutCloserCh <- &timeoutCloser{
		typ:    'r',
		expDur: expDur,
		expTm:  expTm,
		once:   once,
		edr:    epc.Eavesdropper,
		closer: epc}
}

func (epc *EavesdroppedPacketConn) SetWriteTimeout(expired time.Duration, once bool) {
	if epc == nil || epc.PacketConn == nil || expired == 0 {
		return
	}

	just := epc.EnableReadWriteTime()

	t := epc.WriteTime()
	expDur := int64(expired / time.Millisecond)
	expTm := t + expDur
	now := NowUnixMilli()

	if !just && now >= expTm {
		epc.Close()
		return
	}

	if once && epc.PacketConn.SetWriteDeadline(time.UnixMilli(now).Add(5*time.Second)) == nil {
		return
	}

	startTimeoutClosers()
	timeoutCloserCh <- &timeoutCloser{
		typ:    'w',
		expDur: expDur,
		expTm:  expTm,
		once:   once,
		edr:    epc.Eavesdropper,
		closer: epc}
}*/
