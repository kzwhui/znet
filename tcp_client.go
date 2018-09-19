package znet

import "time"

type TCPClient struct {
	connectionPool *ConnectionPool
}

func NewTCPClient(poolnum int, ipPort string, zkname string) *TCPClient {
	return &TCPClient{
		connectionPool: NewConnectionPool(poolnum, ipPort, zkname),
	}
}

func (c *TCPClient) DoRequest(request PacketInterface, response PacketInterface, timeout time.Duration) (err error) {
	conn, err := c.connectionPool.Pop()
	if err != nil {
		return
	}

	defer func() {
		c.connectionPool.Push(conn)
	}()

	if err = conn.Send(request, timeout/2); err != nil {
		conn.Close()
		conn = nil
		return
	}

	if err = conn.Recv(response, timeout/2); err != nil {
		conn.Close()
		conn = nil
		return
	}

	return
}
