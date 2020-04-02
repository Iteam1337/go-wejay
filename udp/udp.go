package udp

import (
	"log"
	"net"
)

func Connect(addr string) (req Req) {
	conn, err := net.Dial("udp4", addr)
	if err != nil {
		log.Fatal(err)
	}

	var acc = make(chan bool, 1)
	go ping(conn, acc)
	if ok := <-acc; !ok {
		log.Fatal("could not send initial ping")
	}

	req.addr = addr
	return
}
