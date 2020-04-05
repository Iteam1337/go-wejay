package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	auth routeAuth
	main routeMain
	room routeRoom
}

type routePath string

var router = route{}

const (
	routePathBase      routePath = "/"
	routePathEmpty     routePath = "/empty"
	routePathLeaveRoom routePath = "/room/leave"
)

func redirect(w http.ResponseWriter, r *http.Request, path interface{}) {
	http.Redirect(w, r, fmt.Sprintf("//%s%s", r.Host, path), 307)
}

func serverListen() {
	r := mux.NewRouter()

	r.HandleFunc("/callback", router.auth.Callback).Methods("GET")
	r.HandleFunc("/sign-out", router.auth.SignOut).Methods("GET")

	r.HandleFunc("/rooms", router.room.Query).Methods("GET")
	r.HandleFunc("/room", router.room.Join).Methods("GET")
	r.HandleFunc("/room/leave", router.room.Leave).Methods("GET")
	r.HandleFunc("/room/{room}", router.room.View).Methods("GET")

	r.HandleFunc("/empty", router.main.Empty).Methods("GET")
	r.HandleFunc("/", router.main.Root).Methods("GET")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", r)

	log.Printf("Listen on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
