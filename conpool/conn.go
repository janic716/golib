package conpool

import (
	"net"
	"time"
	"sync/atomic"
)

var (
	connectId = 0
)

type NewConnFunc func() (net.Conn, error)

type Heartbeat func()

type Connect struct {
	conn       net.Conn
	Borrow     bool
	Close      bool
	ActiveTime int64
	id         int64
	owner      *connSet
	heartbeat  Heartbeat
	usedTimes  int
}

func newConnect(newConnFunc NewConnFunc, heartbeat Heartbeat, owner *connSet, conf Conf) (connect *Connect, err error) {
	connect = &Connect{Borrow:false, Close:false, ActiveTime:0, id: atomic.AddInt64(&connectId, 1), owner:owner, heartbeat:heartbeat}
	conn, err := newConnFunc()
	if err == nil {
		connect.conn = conn
		connect.conn.SetReadDeadline(conf.ReadDeadLine)
		connect.conn.SetWriteDeadline(conf.WriteDeadline)
	}
	return
}

func (c *Connect)updateActiveTime() {
	c.ActiveTime = time.Now().UnixNano() / 1000
	c.conn.Close()
}

func (c *Connect)Read(b []byte) (n int, err error) {
	n, err = c.conn.Read(b)
	if err == nil {
		c.updateActiveTime()
	}
	return
}

func (c *Connect)Write(b []byte) (n int, err error) {
	n, err = c.conn.Write(b)
	if err == nil {
		c.updateActiveTime()
	}
	return
}

func (c *Connect)Closed() error {
	return c.owner.putConn(c)
}

func (c *Connect)LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Connect)RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connect)SetDeadline(t time.Time) error {
	return c.SetDeadline(t)
}

func (c *Connect)SetReadDeadline(t time.Time) error {
	return c.SetReadDeadline(t)
}

func (c *Connect)SetWriteDeadline(t time.Time) error {
	return c.SetWriteDeadline(t)
}

type Conf struct {
	InitConnNum   int
	MaxConnNum    int
	ReadDeadLine  time.Time
	WriteDeadline time.Time
	Block         bool
	BlockTime     time.Time
}

