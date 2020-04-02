package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	auth routeAuth
	main routeMain
	room routeRoom
}

var router = route{}

func serverListen() {
	r := mux.NewRouter()

	r.HandleFunc("/callback", router.auth.callback).Methods("GET")
	r.HandleFunc("/sign-out", router.auth.signOut).Methods("GET")
	r.HandleFunc("/new-auth", router.auth.newAuth).Methods("GET")

	r.HandleFunc("/rooms", router.room.rooms).Methods("GET")
	r.HandleFunc("/room", router.room.join).Methods("GET")
	r.HandleFunc("/room/leave", router.room.leave).Methods("GET")

	r.HandleFunc("/profile", router.main.profile).Methods("GET")
	r.HandleFunc("/", router.main.root).Methods("GET")

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", r)

	log.Printf("Listen on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
