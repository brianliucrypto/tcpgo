package tcpgo

import (
	"fmt"

	"github.com/brianliucrypto/tcpgo/iface"
)

type MessageHandler struct {
	routers  map[uint32]iface.IRouter
	poolSize uint32

	messageChan []chan iface.IRequest

	exitChan chan struct{}
}

func NewMessageHandler() iface.IMessageHandler {
	return &MessageHandler{
		routers:     make(map[uint32]iface.IRouter),
		poolSize:    4,
		messageChan: make([]chan iface.IRequest, 4),
		exitChan:    make(chan struct{}),
	}
}

func (m *MessageHandler) Start() {
	for i := 0; i < int(m.poolSize); i++ {
		go m.receiveMessage(i)
	}
}

func (m *MessageHandler) Stop() {
	m.exitChan <- struct{}{}
	close(m.exitChan)
}

func (m *MessageHandler) AddRouter(msgID uint32, router iface.IRouter) {
	m.routers[msgID] = router
}

func (m *MessageHandler) HandleMessage(message iface.IRequest) {
	router, ok := m.routers[message.GetMessage().GetMsgId()]
	if !ok {
		fmt.Println("router not found")
		return
	}

	router.PreHandle(message)
	router.Handle(message)
	router.PostHandle(message)
}

func (m *MessageHandler) SendMessage2Queue(message iface.IRequest) {
	fmt.Println("receiveMessage start")
	index := message.GetConnection().GetConnID() % m.poolSize
	m.messageChan[index] <- message
}

func (m *MessageHandler) receiveMessage(index int) {
	defer fmt.Println("receiveMessage exit, index:", index)
	fmt.Println("receiveMessage start, index:", index)
	for {
		select {
		case message := <-m.messageChan[index]:
			m.HandleMessage(message)

		case <-m.exitChan:
			return
		}
	}
}
