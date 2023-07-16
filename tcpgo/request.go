package tcpgo

import "github.com/brianliucrypto/tcpgo/iface"

type Request struct {
	Conn    *Connection
	Message *Message
}

func (r *Request) GetConnection() iface.IConnection {
	return r.Conn
}

func (r *Request) GetData() iface.IMessage {
	return r.Message
}
