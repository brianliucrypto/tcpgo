package tcpgo

import (
	"github.com/brianliucrypto/tcpgo/iface"
	"github.com/brianliucrypto/tcpgo/tlog"
)

type MessageHandler struct {
	routers  map[uint32]iface.IRouter
	poolSize uint32

	messageChan []chan iface.IRequest

	exitChan chan struct{}
}

func NewMessageHandler(poolSize uint32) iface.IMessageHandler {
	return &MessageHandler{
		routers:     make(map[uint32]iface.IRouter),
		poolSize:    poolSize,
		messageChan: make([]chan iface.IRequest, 4),
		exitChan:    make(chan struct{}),
	}
}

func (m *MessageHandler) Start() {
	for i := 0; i < int(m.poolSize); i++ {
		m.messageChan[i] = make(chan iface.IRequest, 1024)
		go m.receiveMessage(i, m.messageChan[i])
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
		tlog.Info("router not found")
		return
	}

	router.PreHandle(message)
	router.Handle(message)
	router.PostHandle(message)
	tlog.Info("receiveMessage end")
}

func (m *MessageHandler) SendMessage2Queue(message iface.IRequest) {
	tlog.Info("receiveMessage start")
	index := message.GetConnection().GetConnID() % m.poolSize
	m.messageChan[index] <- message
}

func (m *MessageHandler) receiveMessage(index int, messageChan chan iface.IRequest) {
	defer tlog.Info("receiveMessage exit, index:", index)
	tlog.Info("receiveMessage start, index:", index)
	for {
		select {
		case message := <-messageChan:
			m.HandleMessage(message)

		case <-m.exitChan:
			return
		}
	}
}
