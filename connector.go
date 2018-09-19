package znet

import (
	"bytes"
	"errors"
	"net"
	"time"
)

type Connector struct {
	conn *net.TCPConn
}

func NewConnector(ipPort string) (connector *Connector, err error) {
	connector = &Connector{}
	err = connector.Connect(ipPort, 3*time.Second)
	return
}

func (c *Connector) Connect(ipPort string, timeout time.Duration) (err error) {
	var commConn net.Conn
	if commConn, err = net.DialTimeout("tcp", ipPort, timeout); err != nil {
		return
	}

	var ok bool
	if c.conn, ok = commConn.(*net.TCPConn); !ok {
		err = errors.New("connect timeout")
		return
	}

	return
}

func (c *Connector) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
	c.conn = nil
}

// 同步
func (c *Connector) Send(packet PacketInterface, timeout time.Duration) (err error) {
	if c.conn == nil {
		err = errors.New("no connection")
		return
	}

	c.conn.SetWriteDeadline(time.Now().Add(timeout))
	stream, _ := packet.Marshal()
	//log.Printf("send packet hex: %v\n", hex.EncodeToString(stream))
	var n int
	for total, start := len(stream), 0; start < total; {
		if n, err = c.conn.Write(stream[start:]); err == nil {
			start = start + n
		} else {
			c.Close()
			err = errors.New("connection error")
			return
		}
	}

	return
}

// 同步
func (c *Connector) Recv(packet PacketInterface, timeout time.Duration) (err error) {
	if c.conn == nil {
		err = errors.New("no connection")
		return
	}

	var packetBuf bytes.Buffer
	recvBuf := make([]byte, 1024)
	var n int

	for {
		// 设定读取超时
		c.conn.SetReadDeadline(time.Now().Add(timeout))
		n, err = c.conn.Read(recvBuf[0:])
		if err != nil {
			return
		}

		if n <= 0 {
			c.Close()
			err = errors.New("connection error")
			return
		}

		packetBuf.Write(recvBuf[0:n])
		extra := packet.Check(packetBuf.Bytes())
		if extra < 0 {
			continue
		}

		if err = packet.Unmarshal(packetBuf.Next(packetBuf.Len() - int(extra))); err != nil {
			c.Close()
			return
		} else {
			break
		}
	}

	return
}
