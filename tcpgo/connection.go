package tcpgo

import (
	"fmt"
	"net"
)

type Connection struct {
	Conn   net.Conn
	ConnID uint32

	msgChan chan []byte

	isClose   bool
	ExitChain chan struct{}
}

func NewConneciton(conn net.Conn, connID uint32) *Connection {
	return &Connection{
		Conn:      conn,
		ConnID:    connID,
		msgChan:   make(chan []byte),
		isClose:   false,
		ExitChain: make(chan struct{}),
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
	fmt.Printf("connid:%v stop", c.ConnID)

	c.isClose = true
	c.Conn.Close()

	c.ExitChain <- struct{}{}
	close(c.ExitChain)
	close(c.msgChan)
}

func (c *Connection) StartRead() {
	fmt.Printf("connid:%v reader routine is running", c.ConnID)
	defer c.Stop()

	for {
		readBuf := make([]byte, 1024)
		cnt, err := c.Conn.Read(readBuf)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Print("receive msg", string(readBuf[:cnt]))
		c.msgChan <- readBuf[:cnt]
	}

}

func (c *Connection) StartWrite() {
	fmt.Printf("connid:%v writer routine is running", c.ConnID)
	for {
		select {
		case data, ok := <-c.msgChan:
			if !ok {
				break
			}

			_, err := c.Conn.Write(data)
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
