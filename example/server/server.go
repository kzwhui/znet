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
	// err = errors.New("asdfasf")
	return
}

func (m *MsgCallBack) OnMessage(c *znet.TcpConnection, p znet.PacketInterface) (err error) {
	rawPacket := p.(*znet.Packet)
	fmt.Println("recv: ", string(rawPacket.Data))
	rawPacket.Data = []byte("amazing!!!")
	c.AsyncSend(p, 0)
	fmt.Println("reply: ", string(rawPacket.Data))
	return
}

func (m *MsgCallBack) OnClose(c *znet.TcpConnection) {
	log.Println("onclose: ", c.GetRawConn().RemoteAddr())
}

func main() {
	server := znet.NewTcpServer(&MsgCallBack{})
	if err := server.Start("0.0.0.0:2345"); err != nil {
		fmt.Println(err)
		return
	}

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	server.Stop()
}
