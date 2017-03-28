package smap

import (
	"runtime"
	"fmt"
	"syscall"
)

type Smap struct {
	chunks    []*safeMap
	chunksNum int
}

func NewSmap() *Smap {
	cpuNum := runtime.NumCPU()
	if cpuNum < 1 {
		cpuNum = 1
	}
	chunksNum := 1
	cpuNum--
	for cpuNum > 0{
		cpuNum = cpuNum >> 1
		chunksNum = chunksNum << 1
	}
	chunksNum-- //确保 chunksnum 为 2^n - 1
	chunks := make([]*safeMap, chunksNum)
	for i := 0; i < chunksNum; i++ {
		chunks[i] = NewSafeMap()
	}
	return &Smap{chunks:chunks, chunksNum:chunksNum}
}

func (s *Smap) Get(key interface{}) interface{} {
	syscall.Syscall()
	keyStr := fmt.Printf("%v", key)
}
