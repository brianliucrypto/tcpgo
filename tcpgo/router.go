package tcpgo

import "github.com/brianliucrypto/tcpgo/iface"

type Router struct {
}

func (r *Router) PreHandle(Request iface.IRequest) {
}

func (r *Router) Handle(Request iface.IRequest) {
}

func (r *Router) PostHandle(Request iface.IRequest) {
}
