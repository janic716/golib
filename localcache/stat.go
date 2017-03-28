package localcache

import (
	"sync/atomic"
	"time"
)

type stat struct {
	hit      int64
	get      int64
	set      int64
	del      int64
	gt10     int64 //[10, 100) micro second
	gt100    int64 //[100, 1000)
	gt1000   int64 //[1000, 10000)
	gt10000  int64 //[10000, 100000)
	gt100000 int64 //[1000000, ...)
}

func (s *stat) incrHit(cnt int64) {
	atomic.AddInt64(&s.hit, cnt)
}

func (s *stat) incrGet(cnt int64) {
	atomic.AddInt64(&s.get, cnt)
}

func (s *stat) incrSet(cnt int64) {
	atomic.AddInt64(&s.set, cnt)
}

func (s *stat) incrDel(cnt int64) {
	atomic.AddInt64(&s.del, cnt)
}

func (s *stat) cost(start int64) {
	cost := (time.Now().UnixNano() - start) / 1000
	if cost < 10 {

	} else if cost < 100 {
		atomic.AddInt64(&s.gt10, 1)
	} else if cost < 1000 {
		atomic.AddInt64(&s.gt100, 1)
	} else if cost < 10000 {
		atomic.AddInt64(&s.gt1000, 1)
	} else if cost < 100000 {
		atomic.AddInt64(&s.gt10000, 1)
	} else {
		atomic.AddInt64(&s.gt100000, 1)

	}
}

func (s *stat) reset() {
	atomic.StoreInt64(&s.hit, 0)
	atomic.StoreInt64(&s.get, 0)
	atomic.StoreInt64(&s.set, 0)
	atomic.StoreInt64(&s.del, 0)
	atomic.StoreInt64(&s.gt10, 0)
	atomic.StoreInt64(&s.gt100, 0)
	atomic.StoreInt64(&s.gt1000, 0)
	atomic.StoreInt64(&s.gt10000, 0)
	atomic.StoreInt64(&s.gt100000, 0)
}
