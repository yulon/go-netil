package netil

const kb4 = 4 * 1024
const kb4u64 = uint64(kb4)

var bufferPool = NewLimitedPool(0, func() any {
	return make([]byte, kb4)
})

func LimitBuffersSize(sz uint64) {
	bufferPool.Limit(sz / kb4u64)
}

func AllocBuffer() []byte {
	return bufferPool.Alloc().([]byte)
}

func RecycleBuffer(b []byte) {
	bufferPool.Recycle(b)
}
