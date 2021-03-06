package core

import (
	"sync"
)

const BufSize = 2 * 1024

var pool = sync.Pool{
	New: func() interface{} {
		return make([]byte, BufSize)
	},
}

func NewBytes(size int) []byte {
	if size <= BufSize {
		return pool.Get().([]byte)
	} else {
		return make([]byte, size)
	}
}

func FreeBytes(b []byte) {
	if len(b) <= BufSize {
		pool.Put(b)
	}
}
