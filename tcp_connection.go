package znet

import (
	"bytes"
	"errors"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	SendQueue uint = 20
	RecvQueue uint = 20
)

type TcpConnection struct {
	callback  ConnCallBackInterface // 业务回调函数
	conn      *net.TCPConn          // 由外界写入
	isClose   int32                 // 1--断开 2--连接
	closeOnce sync.Once             // 一个连接关闭，只能关闭一次
	sendChan  chan PacketInterface  // 发送队列
	recvChan  chan PacketInterface  // 接受队列
	closeChan chan struct{}         // 关闭信号
	exitChan  chan struct{}         // 进程退出信号
	wg        *sync.WaitGroup       // 等待所有线程
}

func NewTcpConn(conn *net.TCPConn, callback ConnCallBackInterface, wg *sync.WaitGroup, exitChan chan struct{}) *TcpConnection {
	return &TcpConnection{
		callback:  callback,
		conn:      conn,
		isClose:   0,
		sendChan:  make(chan PacketInterface, SendQueue),
		recvChan:  make(chan PacketInterface, RecvQueue),
		closeChan: make(chan struct{}),
		exitChan:  exitChan,
		wg:        wg,
	}
}

func (c *TcpConnection) StartLoop() {
	if err := c.callback.OnConnect(c); err != nil {
		log.Println("on connect fail")
		c.Close()
		return
	}

	go c.sendLoop()
	go c.recvLoop()
	go c.handleRecv()
}

func (c *TcpConnection) AsyncSend(packet PacketInterface, timeout time.Duration) (err error) {
	if c.IsClose() {
		return
	}

	if timeout == 0 {
		select {
		case c.sendChan <- packet:
			return
		default:
			err = errors.New("send block")
			return
		}
	} else {
		select {
		case c.sendChan <- packet:
			return
		case <-c.closeChan:
			return
		case <-c.exitChan:
			return
		case <-time.After(timeout):
			err = errors.New("send block")
			return
		}
	}

	return
}

func (c *TcpConnection) Close() {
	c.closeOnce.Do(func() {
		c.callback.OnClose(c)
		atomic.StoreInt32(&c.isClose, 1)
		close(c.closeChan)
		close(c.sendChan)
		close(c.recvChan)
		c.conn.Close()
		log.Println("connection close")
	})
}

func (c *TcpConnection) IsClose() bool {
	return atomic.LoadInt32(&c.isClose) == 1
}

func (c *TcpConnection) GetRawConn() *net.TCPConn {
	return c.conn
}

func (c *TcpConnection) sendLoop() {
	c.wg.Add(1)
	defer c.wg.Done()

	for {
		select {
		case packet, ok := <-c.sendChan:
			if c.IsClose() {
				return
			}

			if !ok {
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(time.Duration(300) * time.Millisecond))
			stream, _ := packet.Marshal()
			for total, start := len(stream), 0; start < total; {
				if n, err := c.conn.Write(stream[start:]); err == nil {
					start = start + n
				} else {
					return
				}
			}

		case <-c.closeChan:
			return

		case <-c.exitChan:
			return
		}
	}
}

func (c *TcpConnection) handleRecv() {
	c.wg.Add(1)
	defer c.wg.Done()

	for {
		select {
		case packet, ok := <-c.recvChan:
			if c.IsClose() {
				return
			}

			if !ok {
				return
			}
			c.callback.OnMessage(c, packet)

		case <-c.closeChan:
			return

		case <-c.exitChan:
			return
		}
	}
}

func (c *TcpConnection) recvLoop() {
	c.wg.Add(1)
	defer c.wg.Done()

	var packetBuf bytes.Buffer
	recvBuf := make([]byte, 1024)
	packetTemp := Packet{}
	for {
		// 检查下是否有关闭信号
		select {
		case <-c.closeChan:
			return

		case <-c.exitChan:
			return

		default:
		}

		// 设定读取超时
		c.conn.SetReadDeadline(time.Now().Add(time.Duration(300) * time.Millisecond))
		n, err := c.conn.Read(recvBuf[0:])
		if err != nil {
			// 若为超时，则不算错误
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			} else {
				c.Close()
				return
			}
		}

		// 若超时等待回来发现，连接已断
		if c.IsClose() {
			return
		}

		if n <= 0 {
			c.Close()
			return
		}

		packetBuf.Write(recvBuf[0:n])
		extra := packetTemp.Check(packetBuf.Bytes())
		if extra < 0 {
			continue
		}

		packet := &Packet{}
		if err = packet.Unmarshal(packetBuf.Next(packetBuf.Len() - int(extra))); err != nil {
			log.Println("fail to parse packet")
			c.Close()
			return
		}

		c.recvChan <- packet
	}
}
