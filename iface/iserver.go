package iface

type IServer interface {
	// Start a server
	Start()

	// Stop a server
	Stop()

	// Run a server
	Serve()
}
