package localcache

import (
	"sync"
	"hash/crc32"
	"time"
	"sync/atomic"
	"math/rand"
	"gitlab.niceprivate.com/arch/golib/log"
	"fmt"
	"sort"
	"bytes"
)

const (
	max_entry_num = 10000000
	chunk_num = 100
)

type entry struct {
	key      string
	value    interface{}
	expireAt int64
}

type LocalCache struct {
	maxEntryNum   int
	entryMapArray []map[string]entry
	mutexArray    []*sync.RWMutex
	chunkNum      uint32
	stat
}

var (
	cache *LocalCache
)

func init() {
	if cache == nil {
		cache = newCache()
	}
}

func newCache() (*LocalCache) {
	cache := &LocalCache{}
	cache.chunkNum = chunk_num
	cache.maxEntryNum = max_entry_num
	cache.entryMapArray = make([]map[string]entry, cache.chunkNum)
	cache.mutexArray = make([]*sync.RWMutex, cache.chunkNum)
	for i := range cache.entryMapArray {
		cache.entryMapArray[i] = make(map[string]entry)
		cache.mutexArray[i] = new(sync.RWMutex)
	}

	go func() {
		errChan := make(chan bool, 1)
		loop:
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Error("localcace task err:", err)
					errChan <- true
				}
			}()
			clearNum := max_entry_num / chunk_num
			for {
				select {
				case <-time.Tick(15 * time.Second):
					entryCount := cache.getEntryCount()
					if entryCount > cache.maxEntryNum {
						cache.randomDel(clearNum * 2)
					}
					cache.clearExpiredEntryTask(clearNum)
				}
			}
		}()
		<-errChan
		goto loop
	}()
	return cache
}

func Get(key string) (val interface{}, err error) {
	currentTime := time.Now().UnixNano()
	atomic.AddInt64(&cache.get, 1)
	entryMap, mu := cache.getEntryMapAndMutexByKey(key)
	mu.RLock()
	if e, found := entryMap[key]; found {
		if currentTime > e.expireAt {
			err = ErrNotFound
		} else {
			val = e.value
			cache.stat.incrHit(1)
		}
	} else {
		err = ErrNotFound
	}
	mu.RUnlock()
	cache.stat.cost(currentTime)
	return
}

func (cache *LocalCache)randomDel(num int) {
	numPerChunk := num / chunk_num
	total := 0
	for i, entryMap := range cache.entryMapArray {
		mu := cache.mutexArray[i]
		flag := 0
		mu.Lock()
		remain := int32(len(entryMap) - numPerChunk)
		if remain > 0 {
			delStartIndex := int(rand.Int31n(int32(len(entryMap) - numPerChunk)))
			for k := range entryMap {
				if flag >= delStartIndex {
					delete(entryMap, k)
					total++
				} else {
					flag++
				}
			}
		}
		mu.Unlock()
		total += flag
	}
}

func Set(key string, value interface{}, expire int64) {
	currentTime := time.Now().UnixNano()
	entryMap, mu := cache.getEntryMapAndMutexByKey(key)
	mu.Lock()
	entryMap[key] = entry{key: key, value: value, expireAt: currentTime + expire * 1000000000}
	mu.Unlock()
	cache.stat.incrSet(1)
	cache.stat.cost(currentTime)
}

func Del(key string) {
	currentTime := time.Now().UnixNano()
	entryMap, mu := cache.getEntryMapAndMutexByKey(key)
	mu.Lock()
	delete(entryMap, key)
	mu.Unlock()
	cache.stat.incrDel(1)
	cache.stat.cost(currentTime)
}

func Stat() map[string]string {
	stat := make(map[string]string)
	stat["hit_rate"] = fmt.Sprintf("%.3f", float64(cache.hit) / float64(cache.get))
	stat["hit"] = fmt.Sprint(cache.hit)
	stat["get"] = fmt.Sprint(cache.get)
	stat["set"] = fmt.Sprint(cache.set)
	stat["del"] = fmt.Sprint(cache.del)
	stat["gt10"] = fmt.Sprint(cache.gt10)
	stat["gt100"] = fmt.Sprint(cache.gt100)
	stat["gt1000"] = fmt.Sprint(cache.gt1000)
	stat["gt10000"] = fmt.Sprint(cache.gt10000)
	stat["gt100000"] = fmt.Sprint(cache.gt100000)
	stat["entry"] = fmt.Sprint(cache.getEntryCount())
	return stat
}

func StatString() string {
	stat := Stat()
	var keys []string
	for k, _ := range stat {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buffer bytes.Buffer
	for _, k := range keys {
		buffer.WriteString(k + ":" + stat[k] + " ")
	}
	return buffer.String()
}

func ResetStat() {
	cache.stat.reset()
}

func (cache *LocalCache)clearExpiredEntryTask(clearNum int) {
	currentTime := time.Now().UnixNano()
	var (
		entryMap map[string]entry
		mu *sync.RWMutex
	)
	for i := range cache.mutexArray {
		entryMap = cache.entryMapArray[i]
		mu = cache.mutexArray[i]
		flag := 0
		mu.Lock()
		for k, e := range entryMap {
			if currentTime > e.expireAt {
				delete(entryMap, k)
				flag++
				if (flag >= clearNum) {
					break
				}
			}
		}
		mu.Unlock()
	}
	log.Notice("localcache-stat:", StatString())
}

func (cache *LocalCache)getEntryMapAndMutexByKey(key string) (map[string]entry, *sync.RWMutex) {
	chunkIndex := getCrc32(key) % cache.chunkNum
	entryMap := cache.entryMapArray[chunkIndex]
	mu := cache.mutexArray[chunkIndex]
	return entryMap, mu
}

func getCrc32(content string) uint32 {
	return crc32.ChecksumIEEE([]byte(content))
}

func (cache *LocalCache)getEntryCount() int {
	entryCount := 0
	for i, m := range cache.entryMapArray {
		mux := cache.mutexArray[i]
		mux.RLock()
		entryCount += len(m)
		mux.RUnlock()
	}
	return entryCount
}