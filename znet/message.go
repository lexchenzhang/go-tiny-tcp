package znet

type IMessage interface {
	GetID() uint32
	SetID(uint32)
	GetDataLen() uint32
	SetData([]byte)
	GetData() []byte
	SetDataLen(uint32)
}

type Message struct {
	ID      uint32
	DataLen uint32
	Data    []byte
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		ID:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Message) GetID() uint32 {
	return m.ID
}

func (m *Message) SetID(id uint32) {
	m.ID = id
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
