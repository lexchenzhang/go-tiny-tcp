package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/lexchenzhang/go-tiny-tcp/utils"
)

// TLV sequences for Message

type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack(b []byte) (IMessage, error)
}

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	// Datalen uint32 + ID uint32
	return 8
}

func (dp *DataPack) Pack(msg IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	// write MsgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetID()); err != nil {
		return nil, err
	}

	// write dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// write Data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(b []byte) (IMessage, error) {
	dataBuff := bytes.NewReader(b)
	msg := &Message{}

	// read MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	// read dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// read Data
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("dataLen is too big")
	}

	return msg, nil
}
