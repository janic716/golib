package resource

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"math/rand"
	"time"
)

type DBConf struct {
	Host         []string `json:"host"`
	User         string   `json:"user"`
	Passwd       string   `json:"passwd"`
	Database     string   `json:"database"`
	Timeout      int      `json:"timeout"`
	ReadTimeout  int      `json:"read_timeout"`
	WriteTimeout int      `json:"write_timeout"`
	MaxLifeTime  int      `json:"max_life_time"`
	MaxIdleConns int      `json:"max_idle_conns"`
	MaxOpenConns int      `json:"max_open_conns"`
}

var (
	dbs map[string][]*sql.DB
)

func InitDBPool(dbconfs map[string]DBConf) error {
	dbs = make(map[string][]*sql.DB)
	for node, conf := range dbconfs {
		pool := make([]*sql.DB, 0, 10)
		for _, host := range conf.Host {
			connConf := mysql.Config{
				User:         conf.User,
				Passwd:       conf.Passwd,
				Net:          "tcp",
				Addr:         host,
				DBName:       conf.Database,
				Timeout:      time.Duration(conf.Timeout) * time.Millisecond,
				ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Millisecond,
				WriteTimeout: time.Duration(conf.WriteTimeout) * time.Millisecond,
			}
			db, err := sql.Open("mysql", connConf.FormatDSN())
			if err != nil {
				return err
			}
			if err = db.Ping(); err != nil {
				return err
			}
			db.SetConnMaxLifetime(time.Duration(conf.MaxLifeTime) * time.Second)
			db.SetMaxIdleConns(conf.MaxIdleConns)
			db.SetMaxOpenConns(conf.MaxOpenConns)
			pool = append(pool, db)
		}
		dbs[node] = pool
	}
	return nil
}

func GetDB(node string) (db *sql.DB, err error) {
	pool, ok := dbs[node]
	if !ok {
		return nil, errors.New("node doesn't exist")
	}
	if len(pool) == 1 {
		return pool[0], err
	}
	rand.Seed(time.Now().Unix())
	index := rand.Intn(len(pool))
	return pool[index], err
}
