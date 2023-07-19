package iface

type IConnectionManager interface {
	AddConnection(conn IConnection)
	RemoveConnection(conn IConnection)
	GetConnection(connId uint32) (IConnection, error)
	GetConnectionCount() int
	ClearConnection()
}
