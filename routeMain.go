package main

import (
	"log"
	"net/http"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-wejay/tmpl"
)

type routeMain struct{}

func init() {
	router.main = routeMain{}
}

func (route *routeMain) Root(w http.ResponseWriter, r *http.Request) {
	html := `<a href="/new-auth">new auth</a>`
	exists, _, err := exists(r)

	if err != nil || !exists {
		tmpl.Base(w, html)
		return
	}

	redirect(w, r, routePathEmpty)
}

func (route *routeMain) Empty(w http.ResponseWriter, r *http.Request) {
	exists, userID, err := exists(r)

	if err != nil {
		redirect(w, r, routePathBase)
		return
	}

	if !exists {
		redirect(w, r, routePathNewAuth)
		return
	}

	var cb message.UserRoomResponse
	if err := updServer.NewRequest(
		types.IUserRoom,
		&message.UserRoom{UserId: userID},
		&cb,
	); err != nil {
		log.Println(err)
	}

	if cb.Ok && cb.RoomId != "" {
		redirect(w, r, "/room/"+cb.RoomId)
	}

	w.Header().Set("Content-Type", "text/html")

	tmpl.Empty(w)
}
