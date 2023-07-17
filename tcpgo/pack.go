package tcpgo

import (
	"bytes"
	"encoding/binary"

	"github.com/brianliucrypto/tcpgo/iface"
)

type Pack struct {
}

func NewPack() iface.IPack {
	return &Pack{}
}

func (p *Pack) GetHeadLen() uint32 {
	return 8
}

func (p *Pack) Pack(msg iface.IMessage) ([]byte, error) {
	writer := bytes.NewBuffer([]byte{})
	err := binary.Write(writer, binary.BigEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	err = binary.Write(writer, binary.BigEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}

	err = binary.Write(writer, binary.BigEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

func (p *Pack) Unpack(data []byte) (iface.IMessage, error) {
	message := &Message{}
	reader := bytes.NewReader(data)
	err := binary.Read(reader, binary.BigEndian, &message.ID)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &message.Len)
	if err != nil {
		return nil, err
	}

	return message, nil
}
