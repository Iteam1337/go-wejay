package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-wejay/tmpl"
	"github.com/gorilla/mux"
)

type routeRoom struct {
	re *regexp.Regexp
}

func init() {
	router.room = routeRoom{regexp.MustCompile(`[a-z0-9-]`)}
}

func (route routeRoom) parsedRoomName(input []string) (out string) {
	if input == nil || strings.TrimSpace(input[0]) == "" {
		return
	}

	pre := strings.TrimSpace(input[0])
	pre = strings.ReplaceAll(pre, " ", "-")
	log.Println(pre)

	m := route.re.FindAllString(pre, -1)
	if m == nil {
		return
	}

	out = strings.Join(m, "")

	return
}

func (route *routeRoom) joinRoom(roomID string, userID string) (room message.RefRoom, ok bool) {
	var cb message.JoinRoomResponse
	if err := updServer.NewRequest(
		types.IJoinRoom,
		&message.JoinRoom{
			RoomId: roomID,
			UserId: userID,
		},
		&cb,
	); err != nil {
		log.Println(err)
	}

	if cb.Ok && cb.Room.Id == roomID && cb.UserId == userID {
		ok = true
		room = *cb.Room
	}

	return
}

func (route *routeRoom) leaveRoom(userID string) (ok bool) {
	var cb message.UserLeaveRoomResponse

	fmt.Println("1")
	if err := updServer.NewRequest(
		types.IUserLeaveRoom,
		&message.UserLeaveRoom{UserId: userID},
		&cb,
	); err != nil {
		fmt.Println("2")
		log.Println(err)
		return
	}

	fmt.Println("3")
	if cb.Error != "" {
		fmt.Println("4")
		log.Println(cb.Error)
		return
	}

	fmt.Println("5", cb.Ok, cb.UserId, userID)
	if cb.Ok && cb.UserId == userID {
		fmt.Println("6")
		ok = true
	}

	fmt.Println("7")

	return
}

func (route *routeRoom) Join(w http.ResponseWriter, r *http.Request) {
	exists, userID, err := exists(r)
	roomID := route.parsedRoomName(r.URL.Query()["name"])

	if roomID == "" || userID == "" || !exists || err != nil {
		redirect(w, r, routePathBase)
		return
	}

	room, ok := route.joinRoom(roomID, userID)

	if !ok {
		redirect(w, r, routePathBase)
		return
	} else {
		log.Println(room)
	}

	redirect(w, r, "/room/"+roomID)
}

func (route *routeRoom) Leave(w http.ResponseWriter, r *http.Request) {
	exists, id, err := exists(r)

	log.Println("1")
	if err != nil || id == "" || !exists {
		log.Println("2")
		redirect(w, r, routePathBase)
		return
	}

	log.Println("3")
	if ok := route.leaveRoom(id); !ok {
		log.Println("4")
		redirect(w, r, routePathBase)
		return
	}
	log.Println("5")

	redirect(w, r, routePathEmpty)
}

func (route *routeRoom) Query(w http.ResponseWriter, r *http.Request) {

}

func (route *routeRoom) View(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["room"]
	exists, userID, err := exists(r)

	if err != nil || userID == "" || !exists {
		redirect(w, r, routePathBase)
		return
	}

	var cb message.UserRoomResponse
	if err := updServer.NewRequest(
		types.IUserRoom,
		&message.UserRoom{UserId: userID},
		&cb,
	); err != nil {
		log.Println(err)
		redirect(w, r, routePathBase)
		return
	}

	if cb.Error != "" || !cb.Ok || cb.UserId != userID {
		log.Println(cb.Error)
		redirect(w, r, routePathBase)
		return
	}

	if vars != cb.RoomId {
		redirect(w, r, routePathLeaveRoom)
		return
	}

	tmpl.InRoom(w, vars)
}
