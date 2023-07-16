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
}
