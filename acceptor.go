package znet

import (
	"log"
	"net"
	"sync"
	"time"
)

type Acceptor struct {
	callback  ConnCallBackInterface
	tcpAddr   *net.TCPAddr
	listener  *net.TCPListener
	closeOnce sync.Once // 一个连接关闭，只能关闭一次
	closeChan chan struct{}
	exitChan  chan struct{}
	wg        *sync.WaitGroup
}

func NewAcceptor(host string, wg *sync.WaitGroup, callback ConnCallBackInterface, exitChan chan struct{}) (acceptor *Acceptor, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		return
	}

	acceptor = &Acceptor{
		callback:  callback,
		tcpAddr:   tcpAddr,
		closeChan: make(chan struct{}),
		exitChan:  exitChan,
		wg:        wg,
	}

	return
}

func (a *Acceptor) Listen() (err error) {
	a.listener, err = net.ListenTCP("tcp", a.tcpAddr)
	if err != nil {
		return
	}

	log.Printf("listen addr=%v\n", a.listener.Addr())

	a.wg.Add(1)
	go func() {
		defer func() {
			a.Close()
			a.wg.Done()
			log.Println("acceptor close")
		}()

		for {
			select {
			case <-a.closeChan:
				return

			case <-a.exitChan:
				return

			default:
			}

			a.listener.SetDeadline(time.Now().Add(time.Second))
			conn, err := a.listener.AcceptTCP()
			if err != nil {
				continue
			}

			go NewTcpConn(conn, a.callback, a.wg, a.exitChan).StartLoop()
		}
	}()

	return
}

func (a *Acceptor) Close() {
	a.closeOnce.Do(func() {
		close(a.closeChan)
		a.listener.Close()
	})
}
