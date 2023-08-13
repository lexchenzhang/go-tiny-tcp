package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDatePack(t *testing.T) {
	// create client
	go func() {
		// read msg from client and unpack it
		conn, err := net.Dial("tcp", "127.0.0.1:7777")
		if err != nil {
			fmt.Println("client dial err: ", err)
			return
		}
		defer conn.Close()
		dp := NewDataPack()
		// testing tag-length-value
		// pack msg1
		msg1 := &Message{
			ID:      1,
			DataLen: 4,
			Data:    []byte{'z', 'i', 'n', 'x'},
		}
		sendData1, err := dp.Pack(msg1)
		if err != nil {
			fmt.Println("client pack msg1 error", err)
			return
		}
		// pack msg2
		msg2 := &Message{
			ID:      2,
			DataLen: 6,
			Data:    []byte{' ', 't', 'e', 's', 't', '!'},
		}
		sendData2, err := dp.Pack(msg2)
		if err != nil {
			fmt.Println("client pack msg2 error", err)
			return
		}
		// put msg1 and msg2 together and send to server
		sendData := append(sendData1, sendData2...)
		conn.Write(sendData)
	}()
	// create socketTCP server
	l, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("server accept error", err)
		}
		go func(conn net.Conn) {
			// unpack
			dp := NewDataPack()
			for {
				// read info of pack header
				headerData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headerData)
				if err != nil {
					t.Error("read header error", err)
					break
				}
				msgHeader, err := dp.Unpack(headerData)
				if err != nil {
					t.Error("server unpack err", err)
					return
				}
				if msgHeader.GetDataLen() > 0 {
					// read content of pack data
					msg := msgHeader.(*Message)
					msg.Data = make([]byte, msg.GetDataLen())
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data err", err)
						return
					}
					fmt.Println("->Recv MsgID[", msg.GetID(), "] Len[", msg.GetDataLen(), "] Data[", string(msg.GetData()), "]")
				}
			}
		}(conn)
		return
	}
}
