package netil

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync/atomic"
	"time"
)

func Forward(src, dst net.Conn, buf []byte) {
	ForwardTimeout(src, dst, buf, 0, 0)
}

func ForwardTimeout(src, dst net.Conn, buf []byte, idleTimeout time.Duration, halfClosedIdleTimeout time.Duration) {
	if len(buf) == 0 {
		buf = AllocBuffer()
		defer RecycleBuffer(buf)
	}

	ctx := &struct {
		eofCnt atomic.Uintptr
		expTS  atomic.Int64
	}{}

	if idleTimeout > 0 {
		ctx.expTS.Store(NowUnixMilli() + idleTimeout.Milliseconds())
	}

	Go(func() {
		buf2 := AllocBuffer()
		defer RecycleBuffer(buf2)

		blockedHalfClosedIdleTimeout := false
		for {
			expTS := ctx.expTS.Load()
			if expTS > 0 {
				dst.SetReadDeadline(time.UnixMilli(expTS))
			}

			n, re := dst.Read(buf2)
			if n > 0 {
				if idleTimeout > 0 {
					ctx.expTS.Store(NowUnixMilli() + idleTimeout.Milliseconds())
				} else if !blockedHalfClosedIdleTimeout && ctx.eofCnt.Load() == 1 {
					blockedHalfClosedIdleTimeout = true
					ctx.expTS.Store(0)
					dst.SetReadDeadline(time.Time{})
				}

				we := WriteAll(src, buf2[:n])
				if we != nil {
					src.Close()
					dst.Close()
					return
				}
			}

			if re != nil {
				if re == io.EOF {
					if ctx.eofCnt.Add(1) == 1 {
						dstTCPCon, ok := dst.(*net.TCPConn)
						if ok {
							dstTCPCon.CloseRead()

							srcTCPCon, ok := src.(*net.TCPConn)
							if ok {
								srcTCPCon.CloseWrite()
							}
						}

						oldExpTS := ctx.expTS.Load()
						if halfClosedIdleTimeout > 0 {
							newExpTS := NowUnixMilli() + halfClosedIdleTimeout.Milliseconds()
							if newExpTS > oldExpTS && ctx.expTS.CompareAndSwap(oldExpTS, newExpTS) {
								src.SetReadDeadline(time.UnixMilli(newExpTS))
							}
							return
						}
						if oldExpTS > 0 {
							return
						}
					}
				} else if errors.Is(re, os.ErrDeadlineExceeded) {
					expTS := ctx.expTS.Load()
					if expTS > 0 && NowUnixMilli() < expTS {
						continue
					}
				}
				dst.Close()
				src.Close()
				return
			}
		}
	})

	blockedHalfClosedIdleTimeout := false
	for {
		expTS := ctx.expTS.Load()
		if expTS > 0 {
			src.SetReadDeadline(time.UnixMilli(expTS))
		}

		n, re := src.Read(buf)
		if n > 0 {
			if idleTimeout > 0 {
				ctx.expTS.Store(NowUnixMilli() + idleTimeout.Milliseconds())
			} else if !blockedHalfClosedIdleTimeout && ctx.eofCnt.Load() == 1 {
				blockedHalfClosedIdleTimeout = true
				ctx.expTS.Store(0)
				src.SetReadDeadline(time.Time{})
			}

			we := WriteAll(dst, buf[:n])
			if we != nil {
				dst.Close()
				src.Close()
				return
			}
		}

		if re != nil {
			if re == io.EOF {
				if ctx.eofCnt.Add(1) == 1 {
					srcTCPCon, ok := src.(*net.TCPConn)
					if ok {
						srcTCPCon.CloseRead()

						dstTCPCon, ok := dst.(*net.TCPConn)
						if ok {
							dstTCPCon.CloseWrite()
						}
					}

					oldExpTS := ctx.expTS.Load()
					if halfClosedIdleTimeout > 0 {
						newExpTS := NowUnixMilli() + halfClosedIdleTimeout.Milliseconds()
						if newExpTS > oldExpTS && ctx.expTS.CompareAndSwap(oldExpTS, newExpTS) {
							dst.SetReadDeadline(time.UnixMilli(newExpTS))
						}
						return
					}
					if oldExpTS > 0 {
						return
					}
				}
			} else if errors.Is(re, os.ErrDeadlineExceeded) {
				expTS := ctx.expTS.Load()
				if expTS > 0 && NowUnixMilli() < expTS {
					continue
				}
			}
			src.Close()
			dst.Close()
			return
		}
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
	atomic.StoreInt64(&s.ioTime, NowUnixMilli())
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
					now := NowUnixMilli()
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

	halfClosedTime := NowUnixMilli()
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
		buf = AllocBuffer()
		defer RecycleBuffer(buf)
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
				Go(func() {
					defer pc.Close()

					rbuf := AllocBuffer()
					defer RecycleBuffer(rbuf)

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
				})
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
