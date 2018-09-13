package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kzwhui/znet"
)

type MsgCallBack struct {
}

func (m *MsgCallBack) OnConnect(c *znet.TcpConnection) (err error) {
	log.Println("onconnect: ", c.GetRawConn().RemoteAddr())
	return
}

func (m *MsgCallBack) OnMessage(c *znet.TcpConnection, p znet.PacketInterface) (err error) {
	//log.Println(p.(znet.Packet))
	c.AsyncSend(p, 0)
	return
}

func (m *MsgCallBack) OnClose(c *znet.TcpConnection) {
	log.Println("onclose: ", c.GetRawConn().RemoteAddr())
}

func main() {
	server := znet.NewTcpServer(&MsgCallBack{})
	if err := server.Start(":2345"); err != nil {
		fmt.Println(err)
		return
	}

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	server.Stop()
}
