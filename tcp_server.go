package znet

import (
	"log"
	"sync"
)

type TcpServer struct {
	callback ConnCallBackInterface
	acceptor *Acceptor
	exitChan chan struct{}
	wg       *sync.WaitGroup
}

func NewTcpServer(callback ConnCallBackInterface) *TcpServer {
	return &TcpServer{
		callback: callback,
		exitChan: make(chan struct{}),
		wg:       &sync.WaitGroup{},
	}
}

func (s *TcpServer) Start(host string) (err error) {
	s.acceptor, err = NewAcceptor(host, s.wg, s.callback, s.exitChan)
	if err != nil {
		return
	}

	if err = s.acceptor.Listen(); err != nil {
		return
	}

	log.Println("server start")

	return
}

func (s *TcpServer) Stop() {
	close(s.exitChan)
	s.wg.Wait()
}
