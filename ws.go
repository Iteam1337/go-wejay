package main

import (
	"log"

	"github.com/Iteam1337/go-wejay/cookie"
	"github.com/Iteam1337/go-wejay/udp"
	"golang.org/x/net/websocket"
)

func wsListen(ws *websocket.Conn) {
	var err error
	id, err := cookie.GetID(ws.Request())
	if err != nil {
		return
	}

	log.Printf("ws listen for \"%s\"\n", id)

	msg := make(chan []byte)
	udp := udp.Connect(udpHost, false)
	udp.Listen(msg)

	for {
		res := <-msg
		if err = websocket.Message.Send(ws, string(res)); err != nil {
			log.Println("Can't send")
			break
		}
	}
}
