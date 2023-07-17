package iface

type IMessage interface {
	GetMsgId() uint32

	GetMsgLen() uint32

	GetData() []byte
}
