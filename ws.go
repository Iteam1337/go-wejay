package main

import (
	"encoding/json"
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-wejay/cookie"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
)

func wsListen(ws *websocket.Conn) {
	var err error
	id, err := cookie.GetID(ws.Request())
	if err != nil {
		return
	}

	log.Printf("ws listen for \"%s\"\n", id)

	close := make(chan bool, 1)
	msg := make(chan []byte, 4096)
	go updServer.Listen(&msg, id, &close)

	for {
		buf := <-msg

		out := message.ListenResponse{}
		if err = proto.Unmarshal(buf, &out); err != nil {
			log.Println(err)
			continue
		}

		json, err := json.Marshal(out)
		if err != nil {
			log.Println(err)
			continue
		}

		if err = websocket.Message.Send(ws, string(json)); err != nil {
			log.Println("ws can't send")
			break
		}
	}
}
