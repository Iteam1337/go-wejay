package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const html string = `
<html lang="en">
<head>
	<meta charset="utf8">
	<title>wejay</title>
	<style>ul, ul li { list-style: none; margin: none; padding: 0 }</style>
</head>
<body>%s</body>
`

const list string = `
<ul>
	<li><a href="/player/play">Play</a></li>
	<li><a href="/player/pause">Pause</a></li>
	<li><a href="/player/next">Next track</a></li>
	<li><a href="/player/previous">Previous Track</a></li>
	<li><a href="/player/shuffle">Shuffle</a></li>
</ul>
`

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

		if err := user.RunAction(strings.TrimPrefix(r.URL.Path, "/player/")); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, fmt.Sprintf(html, list))
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
