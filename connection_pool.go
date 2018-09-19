package znet

import (
	"errors"
)

type ConnectionPool struct {
	poolChan chan *Connector
	ipPort   string
}

func NewConnectionPool(num int, ipPort string) *ConnectionPool {
	pool := &ConnectionPool{
		poolChan: make(chan *Connector, num),
		ipPort:   ipPort,
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
			conn, err = NewConnector(c.ipPort)
		}

		return

	default:
		err = errors.New("no idle connection")
		return
	}

	return
}
