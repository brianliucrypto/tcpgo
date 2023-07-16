package iface

type IMessage interface {
	GetMsgId() uint32

	GetMsgLen() uint32

	GetData() []byte

	GetHeadLen() uint32

	Pack() ([]byte, error)

	Unpack([]byte) (IMessage, error)
}
