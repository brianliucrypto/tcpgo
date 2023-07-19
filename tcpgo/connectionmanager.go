package tcpgo

import (
	"errors"
	"sync"

	"github.com/brianliucrypto/tcpgo/iface"
)

type ConnectionManager struct {
	connections map[uint32]iface.IConnection
	lock        sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]iface.IConnection),
	}
}

func (c *ConnectionManager) AddConnection(conn iface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.connections[conn.GetConnID()] = conn
}

func (c *ConnectionManager) RemoveConnection(conn iface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.connections, conn.GetConnID())
}

func (c *ConnectionManager) GetConnection(connId uint32) (iface.IConnection, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	conn, ok := c.connections[connId]
	if !ok {
		return nil, errors.New("connection not found")
	}
	return conn, nil
}

func (c *ConnectionManager) GetConnectionCount() int {
	return len(c.connections)
}

func (c *ConnectionManager) ClearConnection() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for connId, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connId)
	}
}
