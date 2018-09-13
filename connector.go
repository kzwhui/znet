package znet

import (
	"bytes"
	"errors"
	"log"
	"net"
	"time"
)

type Connector struct {
	conn *net.TCPConn
}

func NewConnector() *Connector {
	return &Connector{}
}

func (c *Connector) Connect(ip_port string) (err error) {
	serverAddr, err := net.ResolveTCPAddr("tcp", ip_port)
	if err != nil {
		return
	}

	c.conn, err = net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		return
	}

	return
}

func (c *Connector) Close() {
	c.conn.Close()
}

// 同步
func (c *Connector) Send(packet PacketInterface, timeout time.Duration) (err error) {
	c.conn.SetWriteDeadline(time.Now().Add(timeout))
	stream, _ := packet.Marshal()
	var n int
	for total, start := len(stream), 0; start < total; {
		if n, err = c.conn.Write(stream[start:]); err == nil {
			start = start + n
		} else {
			return
		}
	}

	return
}

// 同步
func (c *Connector) Recv(packet PacketInterface, timeout time.Duration) (err error) {
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
			err = errors.New("connection close")
			return
		}

		packetBuf.Write(recvBuf[0:n])
		extra := packet.Check(packetBuf.Bytes())
		if extra < 0 {
			continue
		}

		if err = packet.Unmarshal(packetBuf.Next(packetBuf.Len() - int(extra))); err != nil {
			log.Println("fail to parse packet")
			c.Close()
			return
		} else {
			break
		}
	}

	return
}
