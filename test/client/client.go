package main

import (
	"fmt"
	"time"

	"github.com/kzwhui/znet"
)

func main() {
	host := "127.0.0.1:2345"
	connector := znet.NewConnector()
	if err := connector.Connect(host); err != nil {
		fmt.Println(err)
		return
	}

	request := &znet.Packet{
		Data: []byte("hello world"),
	}
	fmt.Println("request: ", string(request.Data))

	if err := connector.Send(request, time.Second); err != nil {
		fmt.Println(err)
		return
	}

	response := &znet.Packet{}
	if err := connector.Recv(response, time.Second); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("response: ", string(response.Data))

	return
}
