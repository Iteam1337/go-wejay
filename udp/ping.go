package udp

import (
	"log"
	"net"
)

func ping(conn net.Conn, ok chan bool) {
	buff := make([]byte, 3)
	conn.Write([]byte{'p', 0})
	log.Println("sent Ping")
	_, err := conn.Read(buff)
	if err != nil {
		log.Printf("got error: %s\n", err.Error())
		ok <- false
		return
	}

	if buff[0] != 'P' {
		log.Println("pong not received")
		ok <- false
		return
	}

	log.Println("got Pong")
	ok <- true
}
