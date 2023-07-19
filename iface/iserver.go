package iface

type IServer interface {
	// Start a server
	Start()

	// Stop a server
	Stop()

	// Run a server
	Serve()

	// add router
	AddRouter(msgId uint32, router IRouter)

	// get message handler
	GetMessageHandler() IMessageHandler

	// get connection manager
	GetConnectionManager() IConnectionManager

	// set on connection callback
	SetOnConnStart(func(IConnection))

	// set on connection callback
	SetOnConnStop(func(IConnection))

	// call on connection callback
	CallOnConnStart(conn IConnection)

	// call on connection callback
	CallOnConnStop(conn IConnection)
}
