package tcpgo

import (
	"fmt"
	"io"
	"net"

	"github.com/brianliucrypto/tcpgo/iface"
)

type Connection struct {
	Conn   net.Conn
	ConnID uint32

	msgChan chan iface.IMessage

	msgHandler iface.IMessageHandler

	isClose   bool
	ExitChain chan struct{}
}

func NewConneciton(connID uint32, conn net.Conn, msgHandler iface.IMessageHandler) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		msgHandler: msgHandler,
		msgChan:    make(chan iface.IMessage, 1024),
		isClose:    false,
		ExitChain:  make(chan struct{}),
	}
}

func (c *Connection) Start() {
	go c.StartRead()
	go c.StartWrite()
}

func (c *Connection) Stop() {
	if c.isClose {
		return
	}
	fmt.Printf("connid:%v stop\n", c.ConnID)

	c.isClose = true
	c.Conn.Close()

	c.ExitChain <- struct{}{}
	close(c.ExitChain)
	close(c.msgChan)
}

func (c *Connection) StartRead() {
	fmt.Printf("connid:%v reader routine is running\n", c.ConnID)
	defer c.Stop()

	for {
		packer := NewPack()
		readBuf := make([]byte, packer.GetHeadLen())
		cnt, err := io.ReadFull(c.Conn, readBuf)
		if err != nil {
			fmt.Println(err)
			break
		}

		msgHeader, err := packer.Unpack(readBuf[:cnt])
		if err != nil {
			fmt.Println(err)
			break
		}

		if msgHeader.GetMsgLen() > 0 {
			readBuf := make([]byte, msgHeader.GetMsgLen())
			_, err = io.ReadFull(c.Conn, readBuf)
			if err != nil {
				fmt.Println(err)
				break
			}

			msg := msgHeader.(*Message)
			msg.Data = readBuf
		}

		request := &Request{
			Conn:    c,
			Message: msgHeader.(*Message),
		}

		c.msgHandler.SendMessage2Queue(request)
	}

}

func (c *Connection) StartWrite() {
	fmt.Printf("connid:%v writer routine is running\n", c.ConnID)
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
				fmt.Print("write error", err)
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
