package udp

import (
	"log"
	"net"

	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-protobuf-wejay/version"
)

func ping(conn net.Conn, ok chan bool) {
	buff := make([]byte, 3)

	_, writeErr := conn.Write([]byte{byte(types.IPing), byte(version.Version)})

	if writeErr != nil {
		log.Printf("writeErr\t got error: %s\n", writeErr.Error())
		ok <- false
		return
	}

	_, readErr := conn.Read(buff)
	if readErr != nil {
		log.Printf("readErr\t got error: %s\n", readErr.Error())
		ok <- false
		return
	}

	if buff[0] != byte(types.RPong) {
		ok <- false
		return
	}

	ok <- true
}
