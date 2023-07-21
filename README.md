# tcpgo

This is a long link lightweight server framework based on the idea of ​​[zinx](https://github.com/aceld/zinx)

Currently supported protocols are as follows:
* tcp
* websocket(todo)

server example:
```
package main

import (
	"fmt"

	"github.com/brianliucrypto/tcpgo/iface"
	"github.com/brianliucrypto/tcpgo/tcpgo"
)

type PingRouter struct {
	tcpgo.Router
}

func (p *PingRouter) Handle(request iface.IRequest) {
	message := request.GetMessage()
	fmt.Printf("recv messageID:%v, messgeLen:%v, message:%v\n ", message.GetMsgId(), message.GetMsgLen(), message.GetData())

	data := "ping...ping...ping"
	msg := &tcpgo.Message{
		ID:   1,
		Len:  uint32(len(data)),
		Data: []byte(data),
	}

	d, err := tcpgo.NewPack().Pack(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = request.GetConnection().GetConnection().Write(d)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	server := tcpgo.NewServer()
	server.AddRouter(1, &PingRouter{})
	server.SetOnConnStart(func(conn iface.IConnection) {
		fmt.Println("onConnStart...", conn.GetConnection().RemoteAddr().String())
	})

	server.SetOnConnStop(func(conn iface.IConnection) {
		fmt.Println("onConnStop...", conn.GetConnection().RemoteAddr().String())
	})
	server.Serve()
}
```

client example:
```
package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/brianliucrypto/tcpgo/tcpgo"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		packer := tcpgo.NewPack()
		d, err := packer.Pack(&tcpgo.Message{
			ID:   1,
			Len:  5,
			Data: []byte("hello"),
		})
		if err != nil {
			fmt.Println(err)
			break
		}

		time.Sleep(time.Second)
		_, err = conn.Write(d)
		if err != nil {
			fmt.Println(err)
			break
		}

		buf := make([]byte, packer.GetHeadLen())
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			fmt.Println(err)
			break
		}

		msgHeader, err := packer.Unpack(buf)
		if err != nil {
			fmt.Println(err)
			break
		}

		buf = make([]byte, msgHeader.GetMsgLen())
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(string(buf[:]))
	}
}

```