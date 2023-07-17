package iface

import "net"

type IConnection interface {
	// Start
	Start()

	// Stop
	Stop()

	StartRead()

	StartWrite()

	GetConnID() uint32

	GetConnection() net.Conn

	SendMessage(IRequest) error
}
