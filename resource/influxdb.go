package resource

import (
	"errors"
	"github.com/influxdata/influxdb/client/v2"
	"gitlab.niceprivate.com/arch/golib/pool"
	"time"
)

type InfluxDbConf struct {
	Host      string `json:"host"`
	Username  string `json:"username"`
	Password  string `json:"passwd"`
	Database  string `json:"database"`
	Precision string `json:"precision"`
	InitCap   int    `json:"init_cap"`
	MaxIdle   int    `json:"max_idle"`
}

var (
	iconfs map[string]InfluxDbConf
	ipools map[string]pool.Pool
)

func InitInfluxdbPool(confs map[string]InfluxDbConf) (err error) {
	iconfs = confs
	ipools = make(map[string]pool.Pool, len(confs))
	for node, config := range confs {
		host := config.Host
		username := config.Username
		passwd := config.Password
		newf := func() (pool.Conn, error) {
			c, e := client.NewHTTPClient(client.HTTPConfig{
				Addr:     host,
				Username: username,
				Password: passwd,
			})
			return c, e
		}
		p, err := pool.New(config.InitCap, config.MaxIdle, newf)
		if err != nil {
			return err
		}
		ipools[node] = p
	}
	return
}

func WriteInfluxDB(node string, name string, tags map[string]string, fields map[string]interface{}, logTime time.Time) (err error) {
	p, ok := ipools[node]
	if !ok {
		return errors.New("node doesn't exists")
	}
	conf, ok := iconfs[node]
	if !ok {
		return errors.New("node doesn't exists")
	}

	c, err := p.Get()
	if err != nil {
		return err
	}
	pc, ok := c.Conn.(client.Client)
	if !ok {
		return errors.New("type is error")
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  conf.Database,
		Precision: conf.Precision,
	})
	if err != nil {
		return
	}

	pt, err := client.NewPoint(name, tags, fields, logTime)
	if err != nil {
		return
	}
	bp.AddPoint(pt)

	err = pc.Write(bp)
	if err != nil {
		c.MarkUseless()
	}
	return
}
