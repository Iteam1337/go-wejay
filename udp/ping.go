package udp

import (
	"log"
	"net"

	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-protobuf-wejay/version"
)

func ping(conn net.Conn, ok chan bool) {
	buff := make([]byte, 3)

	conn.Write([]byte{byte(types.IPing), byte(version.Version)})
	_, err := conn.Read(buff)
	if err != nil {
		log.Printf("got error: %s\n", err.Error())
		ok <- false
		return
	}

	if buff[0] != byte(types.RPong) {
		ok <- false
		return
	}

	ok <- true
}
