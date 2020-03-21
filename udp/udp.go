package udp

import (
	"log"
	"net"
)

// Connect …
func Connect(addr string, fatalOnError bool) (req Req) {
	log.Println("starting udp client")
	conn, err := net.Dial("udp4", addr)
	if err != nil {
		if !fatalOnError {
			return
		}
		log.Fatal(err)
	}

	var acc = make(chan bool, 1)
	go ping(conn, acc)
	if ok := <-acc; !ok {
		if !fatalOnError {
			return
		}

		log.Fatal("could not send inital ping")
	}

	req.addr = addr
	return
}
