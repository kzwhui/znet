package znet

import (
	"errors"
	"fmt"

	"code.oa.com/kt/nameapi"
)

type ConnectionPool struct {
	poolChan chan *Connector
	ipPort   string
	zkname   string
}

func NewConnectionPool(num int, ipPort string, zkname string) *ConnectionPool {
	pool := &ConnectionPool{
		poolChan: make(chan *Connector, num),
		ipPort:   ipPort,
		zkname:   zkname,
	}

	for i := 0; i < num; i++ {
		pool.poolChan <- nil
	}

	return pool
}

func (c *ConnectionPool) Push(conn *Connector) (err error) {
	if c.poolChan == nil {
		err = errors.New("pool is nil")
		return
	}

	select {
	case c.poolChan <- conn:
		return

	default:
		if conn != nil {
			conn.Close()
		}
		return
	}

	return
}

func (c *ConnectionPool) Pop() (conn *Connector, err error) {
	if c.poolChan == nil {
		err = errors.New("pool is nil")
		return
	}

	select {
	case conn = <-c.poolChan:
		if conn == nil {
			var ip string
			var port int
			if ip, port, err = nameapi.GetHostByKey(c.zkname); err == nil {
				c.ipPort = fmt.Sprintf("%v:%v", ip, port)
			}
			conn, err = NewConnector(c.ipPort)
		}

		return

	default:
		err = errors.New("no idle connection")
		return
	}

	return
}
