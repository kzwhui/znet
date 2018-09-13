package znet

type ConnCallBackInterface interface {
	OnConnect(*TcpConnection) error
	OnMessage(*TcpConnection, PacketInterface) error
	OnClose(*TcpConnection)
}
