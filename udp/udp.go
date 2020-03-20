package udp

import (
	"fmt"
	"log"
	"net"

	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/golang/protobuf/proto"
)

// Req …
type Req struct {
	conn net.Conn
}

// NewRequest …
func (r *Req) NewRequest(m types.MessageType, in proto.Message, out proto.Message) (err error) {
	conn := r.conn
	data, err := proto.Marshal(in)
	if err != nil {
		log.Println(err)
		return
	}
	ver := m.ByteAndVersion()
	_, err = conn.Write(append(ver[:], data[:]...))
	if err != nil {
		log.Println(err)
		return
	}

	buffer := make([]byte, 4096)
	byteLen, err := conn.Read(buffer)
	if err != nil {
		log.Println(err)
		return
	}

	if byteLen < 2 {
		err = fmt.Errorf("response length; expected at least 2, got %d", byteLen)
		return
	}

	if err = proto.Unmarshal(buffer[2:byteLen], out); err != nil {
		log.Println(err)
		return
	}

	return
}

// Connect …
func Connect(addr string) (req Req) {
	log.Println("starting udp client")
	conn, err := net.Dial("udp4", addr)
	if err != nil {
		log.Fatal(err)
	}

	req.conn = conn

	var acc = make(chan bool, 1)
	go ping(conn, acc)
	if ok := <-acc; !ok {
		log.Println("could not send inital ping")
		return
	}

	return
}
