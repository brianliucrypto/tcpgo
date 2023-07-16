package tcpgo

import "github.com/brianliucrypto/tcpgo/iface"

type Request struct {
	Conn    *Connection
	Message *Message
}

func (r *Request) GetConnection() iface.IConnection {
	return r.Conn
}

func (r *Request) GetMessage() iface.IMessage {
	return r.Message
}
