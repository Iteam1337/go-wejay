package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

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
			http.Error(w, "User not found", http.StatusForbidden)
			return
		}

		url := strings.TrimPrefix(r.URL.Path, "/player/")
		if url != "" {
			if err := user.RunAction(ActionFromString(url)); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			http.Redirect(w, r, "//"+r.Host+"/player", 307)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		TmplPlayer(w, user.NowPlaying())
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
