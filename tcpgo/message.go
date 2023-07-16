package tcpgo

import (
	"bytes"
	"encoding/binary"

	"github.com/brianliucrypto/tcpgo/iface"
)

type Message struct {
	ID   uint32
	Len  uint32
	Data []byte
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

func (m *Message) GetHeadLen() uint32 {
	return 8
}

func (m *Message) Pack() ([]byte, error) {
	writer := bytes.NewBuffer([]byte{})
	err := binary.Write(writer, binary.BigEndian, m.ID)
	if err != nil {
		return nil, err
	}

	err = binary.Write(writer, binary.BigEndian, m.Len)
	if err != nil {
		return nil, err
	}

	err = binary.Write(writer, binary.BigEndian, m.Data)
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

func (m *Message) Unpack(data []byte) (iface.IMessage, error) {
	reader := bytes.NewReader(data)
	err := binary.Read(reader, binary.BigEndian, &m.ID)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &m.Len)
	if err != nil {
		return nil, err
	}

	return m, nil
}
