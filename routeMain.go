package main

import (
	"log"
	"net/http"
	"sort"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-wejay/tmpl"
	"github.com/Iteam1337/go-wejay/utils"
	"github.com/google/uuid"
)

type routeMain struct{}

func init() {
	router.main = routeMain{}
}

func (route *routeMain) Root(w http.ResponseWriter, r *http.Request) {
	exists, _, err := exists(r)

	if err == nil && exists {
		redirect(w, r, routePathEmpty)
		return
	}
	var cb message.CallbackURLResponse
	if err := updServer.NewRequest(
		types.ICallbackURL,
		&message.CallbackURL{UserId: uuid.New().String()},
		&cb,
	); err != nil {
		http.Error(w, "Couldn't get callback-url", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.NewAuth(w, cb.Url)
}

func sortBySize(available []*message.RefRoom) utils.PairList {
	pl := make(utils.PairList, len(available))
	for i, room := range available {
		pl[i] = utils.Pair{
			ID:   room.Id,
			Size: int(room.Size),
		}
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func (route *routeMain) Empty(w http.ResponseWriter, r *http.Request) {
	exists, userID, err := exists(r)

	if err != nil {
		redirect(w, r, routePathBase)
		return
	}

	if !exists {
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
	}

	if cb.Ok && cb.RoomId != "" {
		redirect(w, r, "/room/"+cb.RoomId)
	}

	rooms, _ := router.room.Available("")
	w.Header().Set("Content-Type", "text/html")
	tmpl.Empty(w, sortBySize(rooms))
}
