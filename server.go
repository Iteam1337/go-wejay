package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var (
	html string
	list string
)

func init() {
	if bytes, err := ioutil.ReadFile("./static/html/shell.html"); err != nil {
		panic(err)
	} else {
		html = string(bytes)
	}

	if bytes, err := ioutil.ReadFile("./static/html/player.html"); err != nil {
		panic(err)
	} else {
		list = string(bytes)
	}
}

// ServerListen …
func ServerListen() {
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

		if err := user.RunAction(ActionFromString(strings.TrimPrefix(r.URL.Path, "/player/"))); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, list)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fmt.Sprintf(html, `<a href="/new-auth">new auth</a>`))
	})

	http.HandleFunc("/new-auth", func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		url := spotifyAuth.AuthURL(id)

		user := User{}
		users[id] = &user

		w.Header().Set("Content-Type", "text/html")

		fmt.Fprint(w, fmt.Sprintf(html, fmt.Sprintf(`<a href="%s">sign in</a>`, url)))
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
