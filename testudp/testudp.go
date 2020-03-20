package testudp

import (
	"log"
	"net"
	"time"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/golang/protobuf/proto"
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

// Test …
func Test() {
	log.Println("starting tcp client")
	conn, err := net.Dial("udp4", ":8090")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	messageProto := message.Message{Text: "Hello World", Timestamp: time.Now().Unix()}
	data, err := proto.Marshal(&messageProto)
	if err != nil {
		log.Fatal(err)
	}

	var acc = make(chan bool, 1)
	go ping(conn, acc)
	if ok := <-acc; !ok {
		log.Println("could not send inital ping")
		return
	}

	conn.Write(append([]byte{'m', 0}, data...))

	for {
		buffer := make([]byte, 4096)
		length, err := conn.Read(buffer)
		if err != nil {
			log.Printf("got err: %s\n", err.Error())
			return
		}

		if buffer[0] != 'm' {
			continue
		}

		msg := message.Message{}
		if err = proto.Unmarshal(buffer[2:length], &msg); err != nil {
			log.Println(err)
			continue
		}

		log.Printf("%s\n", msg.Text)

		go ping(conn, acc)
		if ok := <-acc; !ok {
			log.Println("could not send")
			break
		}
	}
}
