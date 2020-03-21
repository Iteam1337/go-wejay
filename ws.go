package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-wejay/cookie"
	"github.com/Iteam1337/go-wejay/jsonresponses"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func wsListen(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	id, err := cookie.GetID(r)
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

		if !out.Ok || out.Error != "" {
			log.Println("ws res err", out.Ok, out.Error)
			continue
		}

		json, _ := json.Marshal(
			jsonresponses.ListenResponseFromMessage(out),
		)

		if err = c.WriteMessage(websocket.TextMessage, json); err != nil {
			log.Println("ws can't send")
			break
		}
	}
}
