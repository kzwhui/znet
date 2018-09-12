package znet

import (
	"net"
	"sync"
)

var (
	SendQueue uint = 20
	RecvQueue uint = 20
)

type TcpConnection struct {
	Callback CallInterface // 业务回调函数
	Conn     *net.TCPConn  // 由外界写入

	isClose   bool
	closeOnce sync.Once // 一个连接关闭，只能关闭一次
	sendChan  chan PacketInterface
	recvChan  chan PacketInterface
	closeChan chan struct{}
}

func (c *TcpConnection) StartLoop() {
	c.sendChan = make(chan PacketInterface, SendQueue)
	c.recvChan = make(chan PacketInterface, RecvQueue)
	c.closeChan = make(chan struct{})

	c.Callback.OnConnect(c)

	go c.sendLoop()
	go c.recvLoop()
	go c.handleRecv()
}

func (c *TcpConnection) Send(packet PacketInterface) {
	c.sendChan <- packet
}

func (c *TcpConnection) Close() {

}

func (c *TcpConnection) sendLoop() {
	for {
		select {
		case packet <- c.sendChan:
		case <-c.closeChan:
			return
		}
	}
}

func (c *TcpConnection) handleRecv() {
	for {
		select {
		case packet <- c.recvChan:
			c.Callback.OnMessage(c, packet)

		case <-c.closeChan:
			return
		}
	}
}

func (c *TcpConnection) recvLoop() {
	for {

	}
}
