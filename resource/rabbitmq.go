package resource

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.niceprivate.com/arch/golib/pool"
	"math/rand"
	"net"
	"time"
)

type RabbitmqConf struct {
	Host             []string `json:"host"`
	User             string   `json:"user"`
	Passwd           string   `json:"passwd"`
	VHost            string   `json:"vhost,omitempty"`
	ExchangeName     string   `json:"exchange_name"`
	QueueNum         int      `json:"queue_num"`
	ConnectTimeout   int      `json:"connect_timeout,omitempty"`
	ReadTimeout      int      `json:"read_timeout"`
	RoutingKeyPrefix string   `json:"routing_key_prefix"`
	InitCap          int      `json:"init_cap"`
	MaxIdle          int      `json:"max_idle"`
	Retry            int      `json:"retry"`
}

var (
	confs  map[string]RabbitmqConf
	cpools map[string][]pool.Pool
)

func InitRabbitmqPool(conf map[string]RabbitmqConf) (err error) {
	confs = conf
	cpools = make(map[string][]pool.Pool, len(confs))
	for node, config := range confs {
		cpool := make([]pool.Pool, 0, len(config.Host))
		for _, host := range config.Host {
			h := host
			var url string
			/// host invalid，select index 0
			if len(config.VHost) > 0 {
				url = fmt.Sprintf("amqp://%s:%s@%s/%s", config.User, config.Passwd, h, config.VHost)
			} else {
				url = fmt.Sprintf("amqp://%s:%s@%s", config.User, config.Passwd, h)
			}

			newf := func() (pool.Conn, error) {
				conn, e := amqp.DialConfig(url, amqp.Config{
					Dial: func(network, addr string) (net.Conn, error) {
						return net.DialTimeout("tcp", h, time.Duration(config.ConnectTimeout)*time.Millisecond)
					},
				})
				return conn, e
			}

			p, err := pool.New(config.InitCap, config.MaxIdle, newf)
			if err != nil {
				return err
			}
			cpool = append(cpool, p)
		}
		cpools[node] = cpool
	}
	return
}

func getRabbitmq(node string) (conn *pool.PoolConn, err error) {
	cpool, ok := cpools[node]
	if !ok {
		return conn, errors.New("node doesn't exist")
	}
	if len(cpool) == 0 {
		return conn, errors.New("pool is empty")
	}

	var p pool.Pool
	if len(cpool) == 1 {
		p = cpool[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		p = cpool[rand.Intn(len(cpool))]
	}
	conn, err = p.Get()
	return
}

func publish(node string, conf RabbitmqConf, key int, mandatory, immediate bool, msg amqp.Publishing) (err error) {
	conn, err := getRabbitmq(node)
	if err != nil {
		return
	}
	defer conn.Close()

	c, ok := conn.Conn.(*amqp.Connection)
	if !ok {
		return errors.New("invalid rabbitmq connection")
	}

	ch, err := c.Channel()
	if err != nil {
		conn.MarkUseless()
		return
	}

	rk := fmt.Sprintf("%s%d", conf.RoutingKeyPrefix, key%conf.QueueNum)
	err = ch.Publish(conf.ExchangeName, rk, mandatory, immediate, msg)
	if err != nil {
		conn.MarkUseless()
		return
	}
	ch.Close()
	return
}

func Publish(node string, key int, mandatory, immediate bool, msg amqp.Publishing) (err error) {
	conf, ok := confs[node]
	if !ok {
		return errors.New("node doesn't exist")
	}
	retry := conf.Retry
	if retry <= 0 {
		retry = 2
	}

	for i := 0; i < retry; i++ {
		err = publish(node, conf, key, mandatory, immediate, msg)
		if err == nil {
			break
		}
	}
	return
}
