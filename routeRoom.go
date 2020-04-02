package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-protobuf-wejay/types"
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

func (route *routeRoom) leaveRoom(userId string) (ok bool) {
	var cb message.UserLeaveRoomResponse
	if err := updServer.NewRequest(
		types.IUserLeaveRoom,
		&message.UserLeaveRoom{UserId: userId},
		&cb,
	); err != nil {
		log.Println(err)
	}

	if cb.Error != "" {
		log.Println(cb.Error)
	}

	if cb.Ok && cb.UserId == userId {
		ok = true
	}

	return
}

func (route *routeRoom) Join(w http.ResponseWriter, r *http.Request) {
	exists, userID, err := exists(r)
	roomID := route.parsedRoomName(r.URL.Query()["name"])

	if roomID == "" || userID == "" || !exists || err != nil {
		http.Redirect(w, r, "//"+r.Host+"/", 307)
		return
	}

	room, ok := route.joinRoom(roomID, userID)

	log.Println("join", room, ok)

	if !ok {
		http.Redirect(w, r, "//"+r.Host+"/", 307)
		return
	} else {
		log.Println(room)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, e := w.Write([]byte(`{"path":"/room/` + roomID + `"}`)); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func (route *routeRoom) Leave(w http.ResponseWriter, r *http.Request) {
	exists, id, err := exists(r)

	log.Println("1")

	if err != nil || id == "" || !exists {
		log.Println("2")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	log.Println("3")
	if ok := route.leaveRoom(id); !ok {
		log.Println("4")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	log.Println("5")

	w.Header().Set("Content-Type", "application/json")
	if _, e := w.Write([]byte(`{"path":"/room/leave"}`)); e != nil {
		log.Println("6")
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}

func (route *routeRoom) Query(w http.ResponseWriter, r *http.Request) {

}
