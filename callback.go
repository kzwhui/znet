package znet

type CallInterface interface {
	OnConnect(*TcpConnection)
	OnMessage(*TcpConnection, PacketInterface)
	OnClose(*TcpConnection)
}
