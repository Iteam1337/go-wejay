package main

import (
	"net/http"

	"github.com/Iteam1337/go-wejay/tmpl"
)

type routeMain struct{}

func init() {
	router.main = routeMain{}
}

func (route *routeMain) Root(w http.ResponseWriter, r *http.Request) {
	html := `<a href="/new-auth">new auth</a>`
	exists, _, err := exists(r)

	if err != nil || !exists {
		tmpl.Base(w, html)
		return
	}

	redirect(w, r, routePathProfile)
}

func (route *routeMain) Profile(w http.ResponseWriter, r *http.Request) {
	exists, _, err := exists(r)

	if err != nil {
		redirect(w, r, routePathBase)
		return
	}

	if !exists {
		redirect(w, r, routePathNewAuth)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	tmpl.Profile(w)
}
