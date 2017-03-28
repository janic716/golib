package conpool
import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type connSet struct {
	remainedList  list.List
	usedMap       map[int64]*Connect
	mutexRemained sync.Mutex
	mutexUsed     sync.Mutex
}

func (q *connSet)getConn(block bool, blockTime time.Time) (*Connect, error) {
	q.mutexRemained.Lock()
	defer q.mutexRemained.Unlock()
	for e := q.remainedList.Front(); e != nil; e = e.Next() {
		conn, ok := e.Value.(Connect)
		if ok {
			if !conn.Closed {
				q.remainedList.Remove(e)
				q.mutexUsed.Lock()
				conn.Borrow = true
				q.usedMap[conn.id] = conn
				q.mutexUsed.Unlock()
				return conn, nil
			}
		} else {
			q.remainedList.Remove(e)
		}
	}
	if block {

	}
	return nil, errors.New("no available conn")
}

func (q *connSet) putConn(conn *Connect) error {
	if conn == nil && conn.Closed {
		return errors.New("connect is invalid")
	}
	conn.Borrow = false
	q.mutexUsed.Lock()
	delete(q.usedMap, conn.id)
	q.mutexUsed.Unlock()
	q.mutexRemained.Lock()
	q.remainedList.PushFront(conn)
	q.mutexRemained.Unlock()
	return nil
}

func (q *connSet) removeConn(conn *Connect) {
	if conn == nil {
		return
	}
	q.mutexUsed.Lock()
	q.mutexUsed.Unlock()
	delete(q.usedMap, conn.id)
}

func (q *connSet) removeClosedConn() {
	q.mutexUsed.Lock()
	q.mutexUsed.Unlock()
	for k, v := range q.usedMap {
		if v != nil && !v.Closed {
			delete(q.usedMap, k)
		}
	}
}

func (q *connSet) usedCount() int {
	q.mutexUsed.Lock()
	defer q.mutexUsed.Unlock()
	cnt := 0
	for _, conn := range q.usedMap {
		if !conn.Closed {
			cnt++
		}
	}
	return cnt
}

func (q *connSet) remainedCount() int {
	q.mutexRemained.Lock()
	defer q.mutexRemained.Unlock()
	return q.remainedList.Len()
}

func (q *connSet) len() int {
	return q.usedCount() + q.remainedCount()
}
