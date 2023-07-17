package iface

type IMessageHandler interface {
	Start()
	Stop()
	AddRouter(msgID uint32, router IRouter)
	HandleMessage(message IRequest)
	SendMessage2Queue(message IRequest)
}
