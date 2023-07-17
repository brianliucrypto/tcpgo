package tcpgo

type Message struct {
	ID   uint32
	Len  uint32
	Data []byte
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		ID:   id,
		Len:  uint32(len(data)),
		Data: data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.ID
}

func (m *Message) GetMsgLen() uint32 {
	return m.Len
}

func (m *Message) GetData() []byte {
	return m.Data
}
