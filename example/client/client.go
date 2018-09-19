package main

import (
	"fmt"
	"time"

	"github.com/kzwhui/znet"
)

func main() {
	host := "127.0.0.1:2345"
	client := znet.NewTCPClient(5, host)

	for i := 0; i < 30; i++ {
		request := &znet.Packet{
			Data: []byte("hello world"),
		}
		fmt.Println("request: ", string(request.Data))

		response := &znet.Packet{}
		if err := client.DoRequest(request, response, 3*time.Second); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("response: ", string(response.Data))
		time.Sleep(1 * time.Second)
	}

	return
}
