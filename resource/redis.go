package resource

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"gitlab.niceprivate.com/arch/golib/log"
	"math/rand"
	"time"
)

type RedisConf struct {
	Host         []string `json:"host"`
	Passwd       string   `json:"passwd"`
	ConnTimeout  int      `json:"connect_timeout"` //ms
	ReadTimeout  int      `json:"read_timeout"`    //ms
	WriteTimeout int      `json:"write_timeout"`   // ms
	MaxIdle      int      `json:"max_idle"`
	MaxActive    int      `json:"max_active"`
	IdleTimeout  int      `json:"idle_timeout"` //s
}

var (
	redisPool map[string][]*redis.Pool
)

func InitRedisPool(redisConfs map[string]RedisConf) (err error) {
	redisPool = make(map[string][]*redis.Pool)
	for node, conf := range redisConfs {
		pools := make([]*redis.Pool, 0, 10)
		for _, host := range conf.Host {
			pool := &redis.Pool{
				MaxIdle:     conf.MaxIdle,
				MaxActive:   conf.MaxActive,
				IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
				Dial: func() (redis.Conn, error) {
					c, err := redis.DialTimeout(
						"tcp",
						host,
						time.Duration(conf.ConnTimeout)*time.Millisecond,
						time.Duration(conf.ReadTimeout)*time.Millisecond,
						time.Duration(conf.WriteTimeout)*time.Millisecond)
					if err != nil {
						log.Warning(err)
						return nil, err
					}
					if len(conf.Passwd) > 0 {
						if _, err = c.Do("AUTH", conf.Passwd); err != nil {
							c.Close()
							log.Warning(err)
							return nil, err
						}
					}
					return c, err
				},
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err = c.Do("PING")
					if err != nil {
						log.Warning(err)
					}
					return err
				},
			}
			pools = append(pools, pool)
		}

		redisPool[node] = pools
	}
	return
}

func GetRedis(node string) (c redis.Conn, err error) {
	pools, ok := redisPool[node]
	if !ok {
		return c, errors.New("node doesn't exist")
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(pools))
	return pools[index].Get(), nil
}
