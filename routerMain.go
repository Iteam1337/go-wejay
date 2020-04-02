package main

import (
	"fmt"
	"net/http"

	"github.com/Iteam1337/go-wejay/tmpl"
)

type routerMain struct{}

func init() {
	router.main = routerMain{}
}

func (route *routerMain) root(w http.ResponseWriter, r *http.Request) {
	html := `<a href="/new-auth">new auth</a>`
	exists, _, err := exists(r)

	if err != nil || !exists {
		tmpl.Base(w, html)
		return
	}

	http.Redirect(w, r, "//"+r.Host+"/profile", 307)
}

func (route *routerMain) profile(w http.ResponseWriter, r *http.Request) {
	exists, _, err := exists(r)

	if err != nil {
		http.Redirect(w, r, "//"+r.Host+"/", 307)
		return
	}

	if !exists {
		http.Redirect(w, r, "//"+r.Host+"/new-auth", 307)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	tmpl.Profile(w, fmt.Sprintf("ws://%s/ws", addr))
}
