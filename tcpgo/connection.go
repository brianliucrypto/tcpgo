package tcpgo

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/brianliucrypto/tcpgo/iface"
	"github.com/brianliucrypto/tcpgo/tlog"
)

type Connection struct {
	conn   net.Conn
	connID uint32

	msgChan chan iface.IMessage

	server iface.IServer

	isClose bool

	ctx        context.Context
	cancelFunc context.CancelFunc

	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewConneciton(server iface.IServer, connID, writeCacheSize uint32, conn net.Conn) *Connection {
	return &Connection{
		server:   server,
		conn:     conn,
		connID:   connID,
		msgChan:  make(chan iface.IMessage, writeCacheSize),
		isClose:  false,
		property: make(map[string]interface{}),
	}
}

func (c *Connection) Start() {
	c.ctx, c.cancelFunc = context.WithCancel(context.Background())
	c.server.CallOnConnStart(c)
	c.server.GetConnectionManager().AddConnection(c)

	go c.StartRead()
	go c.StartWrite()

	<-c.ctx.Done()
	c.finalizer()
}

func (c *Connection) Stop() {
	c.cancelFunc()
}

func (c *Connection) finalizer() {
	if c.isClose {
		return
	}
	tlog.Info("connid:%v stop", c.connID)

	c.isClose = true
	c.conn.Close()

	c.server.GetConnectionManager().RemoveConnection(c)
	c.server.CallOnConnStop(c)

	close(c.msgChan)
}

func (c *Connection) StartRead() {
	tlog.Info("connid:%v reader routine is running", c.connID)
	defer c.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			packer := NewPack()
			readBuf := make([]byte, packer.GetHeadLen())
			cnt, err := io.ReadFull(c.conn, readBuf)
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
				_, err = io.ReadFull(c.conn, readBuf)
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
}

func (c *Connection) StartWrite() {
	defer tlog.Info("connid:%v writer routine is stoping", c.connID)
	tlog.Info("connid:%v writer routine is running", c.connID)
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			data, ok := <-c.msgChan
			if !ok {
				break
			}

			d, err := NewPack().Pack(data)
			if err != nil {
				break
			}

			_, err = c.conn.Write(d)
			if err != nil {
				tlog.Info("write error:%v", err)
			}
		}
	}
}

func (c *Connection) GetConnID() uint32 {
	return c.connID

}

func (c *Connection) GetConnection() net.Conn {
	return c.conn
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
