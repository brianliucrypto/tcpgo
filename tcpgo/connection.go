package tcpgo

import (
	"errors"
	"io"
	"net"
	"sync"

	"github.com/brianliucrypto/tcpgo/iface"
	"github.com/brianliucrypto/tcpgo/tlog"
)

type Connection struct {
	Conn   net.Conn
	ConnID uint32

	msgChan chan iface.IMessage

	server iface.IServer

	isClose   bool
	ExitChain chan struct{}

	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewConneciton(server iface.IServer, connID, writeCacheSize uint32, conn net.Conn) *Connection {
	return &Connection{
		server:    server,
		Conn:      conn,
		ConnID:    connID,
		msgChan:   make(chan iface.IMessage, writeCacheSize),
		isClose:   false,
		ExitChain: make(chan struct{}),
		property:  make(map[string]interface{}),
	}
}

func (c *Connection) Start() {
	go c.StartRead()
	go c.StartWrite()

	c.server.GetConnectionManager().AddConnection(c)
	c.server.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	if c.isClose {
		return
	}
	tlog.Info("connid:%v stop\n", c.ConnID)

	c.isClose = true
	c.Conn.Close()

	c.server.GetConnectionManager().RemoveConnection(c)
	c.server.CallOnConnStop(c)

	c.ExitChain <- struct{}{}
	close(c.ExitChain)
	close(c.msgChan)
}

func (c *Connection) StartRead() {
	tlog.Info("connid:%v reader routine is running\n", c.ConnID)
	defer c.Stop()

	for {
		packer := NewPack()
		readBuf := make([]byte, packer.GetHeadLen())
		cnt, err := io.ReadFull(c.Conn, readBuf)
		if err != nil {
			tlog.Info(err.Error())
			break
		}

		msgHeader, err := packer.Unpack(readBuf[:cnt])
		if err != nil {
			tlog.Info(err.Error())
			break
		}

		if msgHeader.GetMsgLen() > 0 {
			readBuf := make([]byte, msgHeader.GetMsgLen())
			_, err = io.ReadFull(c.Conn, readBuf)
			if err != nil {
				tlog.Info(err.Error())
				break
			}

			msg := msgHeader.(*Message)
			msg.Data = readBuf
		}

		request := &Request{
			Conn:    c,
			Message: msgHeader.(*Message),
		}

		c.server.GetMessageHandler().SendMessage2Queue(request)
	}

}

func (c *Connection) StartWrite() {
	tlog.Info("connid:%v writer routine is running\n", c.ConnID)
	for {
		select {
		case data, ok := <-c.msgChan:
			if !ok {
				break
			}

			d, err := NewPack().Pack(data)
			if err != nil {
				break
			}

			_, err = c.Conn.Write(d)
			if err != nil {
				tlog.Info("write error", err)
			}

		case <-c.ExitChain:
			return
		}
	}
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID

}

func (c *Connection) GetConnection() net.Conn {
	return c.Conn
}

func (c *Connection) SendMessage(request iface.IRequest) error {
	c.msgChan <- request.GetMessage()
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	val, ok := c.property[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	return val, nil
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
