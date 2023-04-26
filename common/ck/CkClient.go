package ck

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/gogf/gf/frame/g"
	"strings"
	"time"
)

type CkClient struct {
	conn driver.Conn
}

func NewCkClient(database string) (error, *CkClient) {
	addr := strings.Split(g.Cfg().GetString("ck.addr"), ",")
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: addr,
		Auth: clickhouse.Auth{
			Database: database,
			Username: g.Cfg().GetString("ck.username"),
			Password: g.Cfg().GetString("ck.password"),
		},
		//Debug:           true,
		DialTimeout:     time.Duration(g.Cfg().GetInt("ck.timeout")) * time.Second,
		MaxOpenConns:    g.Cfg().GetInt("ck.maxopenconn"),
		MaxIdleConns:    g.Cfg().GetInt("ck.maxidleconn"),
		ConnMaxLifetime: time.Duration(g.Cfg().GetInt("ck.ConnMaxLifetime")) * time.Minute,
	})

	if err != nil {
		return err, nil
	}

	return nil, &CkClient{conn: conn}
}

func (t *CkClient) GetConn() driver.Conn {
	return t.conn
}
