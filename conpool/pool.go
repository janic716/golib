package conpool
import (
	"math/rand"
	"github.com/anacrolix/sync"
)

type PoolConf struct {
	
}

type ConnPool struct {
	connListTable map[int]connSet
	mu            sync.RWMutex
	len           int64
	initConnNum   int
	maxConnNum    int
	tableLen      int
}

type Pool interface {
	Len() int
	GetConn() (*Connect, error)
	CloseConn()
	SetConnConf()
	ClearAll()
}

func (this *ConnPool)Len() int {
	cnt := 0
	for _, list := range this.connListTable {
		cnt += list.len()
	}
	return cnt
}

func (cp *ConnPool)Remained() int {
	cnt := 0
	for _, list := range cp.connListTable {
		cnt += list.remainedCount()
	}
	return cnt
}

func (cp *ConnPool)GetConn(block bool) (conn Connect, err error) {
	rn := rand.Int() % cp.tableLen
	connList, found := cp.connListTable[rn]
	if found {
		conn, err = connList.getConn()
		if err == nil {
			return
		} else {
			conn, err = cp.GetConn(block)
		}
	} else {
		cp.tableLen = len(cp.connListTable)
	}
}
