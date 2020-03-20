package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	json "github.com/Iteam1337/go-wejay/jsonResponses"
	"github.com/google/uuid"
)

// ServerListen …
func ServerListen() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
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
	})

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

		artists, track := user.NowPlaying()
		var tpl bytes.Buffer
		TmplNowPlaying(&tpl, artists, track)
		TmplPlayer(w, tpl.String())
	})

	http.HandleFunc("/action/", func(w http.ResponseWriter, r *http.Request) {
		id, err := GetIDFromCookie(w, r)
		if err != nil {
			json.NewPlayerResponseErr(w, "missing id")
			return
		}

		user, ok := users[string(id)]
		if !ok {
			json.NewPlayerResponseErr(w, "user not found")
			return
		}

		action := strings.TrimPrefix(r.URL.Path, "/action/")
		if err := user.RunAction(ActionFromString(action)); err != nil {
			json.NewPlayerResponseErr(w, err.Error())
			return
		}

		json.NewPlayerResponse(w, fmt.Sprintf("action %s applied successful", action))
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

	log.Printf("Listen on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
