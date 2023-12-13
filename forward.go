package netil

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"
)

func Forward(src, dst net.Conn, buf []byte) {
	ForwardTimeout(src, dst, buf, 3*time.Minute, 8*time.Second)
}

type forwardContext struct {
	idleEnd    atomic.Int64
	halfClosed atomic.Uintptr
}

func (fc *forwardContext) StartCloser(close func()) {
	go func() {
		for {
			idleEnd := fc.idleEnd.Load()
			if idleEnd == 0 {
				panic(idleEnd)
			}
			delay := idleEnd - time.Now().UnixMilli()
			if delay <= 0 {
				close()
				return
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}()
}

func ForwardTimeout(src, dst net.Conn, buf []byte, idleTimeout time.Duration, halfIdleTimeout time.Duration) {
	if len(buf) == 0 {
		buf = make([]byte, 4096)
	}

	closeAll := func() {
		src.Close()
		dst.Close()
	}

	fc := &forwardContext{}

	if idleTimeout > 0 {
		fc.idleEnd.Store(time.Now().UnixMilli() + idleTimeout.Milliseconds())
		fc.StartCloser(closeAll)
	}

	go func() {
		buf2 := make([]byte, 4096)

		for {
			n, err := dst.Read(buf2)

			if n > 0 {
				if idleTimeout > 0 {
					fc.idleEnd.Store(time.Now().UnixMilli() + idleTimeout.Milliseconds())
				} else if halfIdleTimeout > 0 && fc.halfClosed.Load() == 1 {
					fc.idleEnd.Store(time.Now().UnixMilli() + halfIdleTimeout.Milliseconds())
				}

				we := WriteAll(src, buf2[:n])
				if we != nil && err == nil {
					err = we
				}
			}

			if err == nil {
				continue
			}

			if fc.halfClosed.Swap(1) != 0 {
				closeAll()
				return
			}

			if err == io.EOF {
				dstTCPCon, ok := dst.(*net.TCPConn)
				if ok {
					dstTCPCon.CloseRead()
				}
				srcTCPCon, ok := src.(*net.TCPConn)
				if ok {
					srcTCPCon.CloseWrite()
				}
			}

			fc.idleEnd.Store(time.Now().UnixMilli() + halfIdleTimeout.Milliseconds())
			if idleTimeout == 0 {
				fc.StartCloser(closeAll)
			}
			return
		}
	}()

	for {
		n, err := src.Read(buf)

		if n > 0 {
			if idleTimeout > 0 {
				fc.idleEnd.Store(time.Now().UnixMilli() + idleTimeout.Milliseconds())
			} else if halfIdleTimeout > 0 && fc.halfClosed.Load() == 1 {
				fc.idleEnd.Store(time.Now().UnixMilli() + halfIdleTimeout.Milliseconds())
			}

			we := WriteAll(dst, buf[:n])
			if we != nil && err == nil {
				err = we
			}
		}

		if err == nil {
			continue
		}

		if fc.halfClosed.Swap(1) != 0 {
			closeAll()
			return
		}

		if err == io.EOF {
			srcTCPCon, ok := src.(*net.TCPConn)
			if ok {
				srcTCPCon.CloseRead()
			}
			dstTCPCon, ok := dst.(*net.TCPConn)
			if ok {
				dstTCPCon.CloseWrite()
			}
		}

		fc.idleEnd.Store(time.Now().UnixMilli() + halfIdleTimeout.Milliseconds())
		if idleTimeout == 0 {
			fc.StartCloser(closeAll)
		}
		return
	}
}

/*type forwardHalfState struct {
	isClosed uintptr
	ioTime   int64
}

func (s *forwardHalfState) IsClosed() bool {
	return atomic.LoadUintptr(&s.isClosed) == 1
}

func (s *forwardHalfState) Close() {
	atomic.StoreUintptr(&s.isClosed, 1)
}

func (s *forwardHalfState) IOTime() int64 {
	return atomic.LoadInt64(&s.ioTime)
}

func (s *forwardHalfState) ReadyIO() {
	atomic.StoreInt64(&s.ioTime, time.Now().UnixMilli())
}

type ForwardCloser struct {
	mtx sync.Mutex

	rState forwardHalfState
	wState forwardHalfState

	err error
}

func NewForwardCloser() *ForwardCloser {
	return &ForwardCloser{}
}

func (fc *ForwardCloser) ReadyRead() {
	fc.rState.ReadyIO()
}

func (fc *ForwardCloser) ReadyWrite() {
	fc.wState.ReadyIO()
}

func (fc *ForwardCloser) CloseRead(close func() error) error {
	return fc.closeHalf(&fc.rState, &fc.wState, close)
}

func (fc *ForwardCloser) CloseWrite(close func() error) error {
	return fc.closeHalf(&fc.wState, &fc.rState, close)
}

type closedHalf struct {
	fc                  *ForwardCloser
	halfClosedTime      int64
	otherHalfState      *forwardHalfState
	otherHalfPrevIOTime int64
	close               func() error
}

var halfClosedKillerOnec sync.Once
var halfClosedKillerCh chan *closedHalf

func startHalfClosedKiller() {
	halfClosedKillerOnec.Do(func() {
		halfClosedKillerCh = make(chan *closedHalf, 4096)
		go func() {
			sccs := list.New()
			for {
				select {
				case scc := <-halfClosedKillerCh:
					sccs.PushBack(scc)

				case <-time.Tick(time.Second):
					now := time.Now().UnixMilli()
					for it := sccs.Front(); it != nil; {
						scc := it.Value.(*closedHalf)

						if scc.otherHalfState.IOTime() != scc.otherHalfPrevIOTime || scc.otherHalfState.IsClosed() {
							next := it.Next()
							sccs.Remove(it)
							it = next
							continue
						}

						if now-scc.halfClosedTime < 3000 {
							it = it.Next()
							continue
						}

						scc.fc.mtx.Lock()
						if !scc.otherHalfState.IsClosed() {
							scc.fc.err = scc.close()
							scc.otherHalfState.Close()
						}
						scc.fc.mtx.Unlock()

						next := it.Next()
						sccs.Remove(it)
						it = next
					}
				}
			}
		}()
	})
}

func (fc *ForwardCloser) closeHalf(curHalfState, otherHalfState *forwardHalfState, close func() error) error {
	if curHalfState.IsClosed() {
		if otherHalfState.IsClosed() {
			return fc.err
		}
		return nil
	}

	halfClosedTime := time.Now().UnixMilli()
	otherHalfIOTime := otherHalfState.IOTime()

	fc.mtx.Lock()
	defer fc.mtx.Unlock()

	if otherHalfState.IsClosed() {
		if curHalfState.IsClosed() {
			return fc.err
		}
		fc.err = close()
		curHalfState.Close()
		return fc.err
	} else if curHalfState.IsClosed() {
		return nil
	}
	curHalfState.Close()

	otherHalfIOTime2 := otherHalfState.IOTime()
	if otherHalfIOTime2 > halfClosedTime || otherHalfIOTime2 != otherHalfIOTime {
		return nil
	}

	startHalfClosedKiller()
	halfClosedKillerCh <- &closedHalf{
		fc:                  fc,
		halfClosedTime:      halfClosedTime,
		otherHalfState:      otherHalfState,
		otherHalfPrevIOTime: otherHalfIOTime,
		close:               close,
	}
	return nil
}

func (fc *ForwardCloser) lockAndCloseAll(close func() error) error {
	fc.mtx.Lock()
	defer fc.mtx.Unlock()

	if !fc.rState.IsClosed() || !fc.wState.IsClosed() {
		fc.err = close()
		fc.rState.Close()
		fc.wState.Close()
	}
	return fc.err
}

func (fc *ForwardCloser) CloseAll(close func() error) error {
	if fc.rState.IsClosed() && fc.wState.IsClosed() {
		return fc.err
	}
	return fc.lockAndCloseAll(close)
}*/

type PacketForwardConn interface {
	ReadFromForward(b []byte) (n int, raddr, faddr net.Addr, err error)
	WriteToForward(data []byte, raddr, faddr net.Addr) (int, error)
	LocalAddr() net.Addr
	Close() error
}

func ForwardPacket(pfc PacketForwardConn, buf []byte, listenPubToken string, listenPub func() (net.PacketConn, error), pubPool *PacketConnPool) error {
	defer pfc.Close()

	pfcAddrStr := pfc.LocalAddr().String()

	if pubPool == nil {
		pubPool = NewPacketConnPool(128)
		defer pubPool.Close()
	}

	if len(buf) == 0 {
		buf = make([]byte, 4096)
	}

	for {
		n, dest, faddr, err := pfc.ReadFromForward(buf)
		if err != nil {
			return err
		}

		pubVAddr := fmt.Sprint(faddr.String(), "->", pfcAddrStr)
		for {
			pubCon, err := pubPool.Get(pubVAddr, listenPubToken, func() (net.PacketConn, error) {
				pc, err := ListenPacketOrDirect(listenPub)
				if err != nil {
					return nil, err
				}
				go func() {
					defer pc.Close()

					rbuf := make([]byte, 4096)
					for {
						n, from, err := pc.ReadFrom(rbuf)
						if err != nil {
							return
						}
						_, err = pfc.WriteToForward(rbuf[:n], from, faddr)
						if err != nil {
							return
						}
					}
				}()
				return pc, nil
			})
			if err != nil {
				break
			}

			_, err = pubCon.WriteTo(buf[:n], dest)
			if err == nil {
				break
			}
			pubPool.Del(pubVAddr, listenPubToken)
		}
	}
}
