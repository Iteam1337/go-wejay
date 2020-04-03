package main

import (
	"encoding/json"
	"errors"
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

	if err := updServer.NewRequest(
		types.IUserLeaveRoom,
		&message.UserLeaveRoom{UserId: userID},
		&cb,
	); err != nil {
		log.Println(err)
		return
	}

	if cb.Error != "" {
		log.Println(cb.Error)
		return
	}

	if cb.Ok && cb.UserId == userID {
		ok = true
	}

	return
}

func (route routeRoom) Available(name string) (rooms []*message.RefRoom, err error) {
	var cb message.QueryRoomsResponse
	if e := updServer.NewRequest(
		types.IQueryRooms,
		&message.QueryRooms{Name: name},
		&cb,
	); e != nil {
		err = e
		return
	}

	if !cb.Ok || cb.Error != "" {
		log.Println(cb.Error)
		err = errors.New("could not query")
		return
	}

	rooms = cb.Room

	return
}

func (route *routeRoom) Join(w http.ResponseWriter, r *http.Request) {
	exists, userID, err := exists(r)
	roomID := route.parsedRoomName(r.URL.Query()["name"])

	if roomID == "" || userID == "" || !exists || err != nil {
		redirect(w, r, routePathBase)
		return
	}

	_, ok := route.joinRoom(roomID, userID)
	if !ok {
		redirect(w, r, routePathBase)
		return
	}

	redirect(w, r, "/room/"+roomID)
}

func (route *routeRoom) Leave(w http.ResponseWriter, r *http.Request) {
	exists, id, err := exists(r)

	if err != nil || id == "" || !exists {
		redirect(w, r, routePathBase)
		return
	}

	if ok := route.leaveRoom(id); !ok {
		redirect(w, r, routePathBase)
		return
	}

	redirect(w, r, routePathEmpty)
}

func (route *routeRoom) Query(w http.ResponseWriter, r *http.Request) {
	var name string
	if query := r.URL.Query()["name"]; len(query) > 0 {
		name = strings.TrimSpace(query[0])
	}

	res, _ := route.Available(name)

	json, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(json); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
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
