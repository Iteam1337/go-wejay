package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// PlayerResponse …
type PlayerResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func (p *PlayerResponse) Write(w http.ResponseWriter) {
	json, err := json.Marshal(p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(p)
	fmt.Println(string(json))

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// NewPlayerResponse …
func NewPlayerResponse(w http.ResponseWriter, message string) (res PlayerResponse) {
	res.Ok = true
	res.Message = message
	res.Write(w)
	return
}

// NewPlayerResponseErr …
func NewPlayerResponseErr(w http.ResponseWriter, message string) (res PlayerResponse) {
	res.Ok = false
	res.Message = message
	res.Error = NewCustomError(message)
	res.Write(w)
	return
}

// ServerListen …
func ServerListen() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/callback", completeAuth)

	http.HandleFunc("/player/", func(w http.ResponseWriter, r *http.Request) {
		id, err := GetIDFromCookie(w, r)
		if err != nil {
			return
		}

		user, ok := users[string(id)]
		if !ok {
			http.Redirect(w, r, "//"+r.Host+"/new-auth", 307)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		TmplPlayer(w, user.NowPlaying())
	})

	http.HandleFunc("/action/", func(w http.ResponseWriter, r *http.Request) {
		id, err := GetIDFromCookie(w, r)
		if err != nil {
			NewPlayerResponseErr(w, "missing id")
			return
		}

		user, ok := users[string(id)]
		if !ok {
			NewPlayerResponseErr(w, "user not found")
			return
		}

		action := strings.TrimPrefix(r.URL.Path, "/action/")
		if err := user.RunAction(ActionFromString(action)); err != nil {
			NewPlayerResponseErr(w, "action failed")
			return
		}

		NewPlayerResponse(w, fmt.Sprintf("action %s applied successful", action))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		TmplBase(w, `<a href="/new-auth">new auth</a>`)
	})

	http.HandleFunc("/new-auth", func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		url := spotifyAuth.AuthURL(id)

		user := User{}
		users[id] = &user

		w.Header().Set("Content-Type", "text/html")

		TmplNewAuth(w, url)
	})

	fmt.Printf("Listen on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	id := CreateAndSetCookie(w, r)
	token, err := spotifyAuth.Token(id, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if user, ok := users[id]; ok {
		user.SetClient(token)
		http.Redirect(w, r, "//"+r.Host+"/player", 307)
	} else {
		http.NotFound(w, r)
	}
}
