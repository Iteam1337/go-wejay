package main

import (
	"html/template"
	"io"
)

var (
	tmpl *template.Template
)

type tmplBase struct {
	HTML template.HTML
}

type tmplNewAuth struct {
	SignIn string
}

type tmplPlayer struct {
	NowPlaying template.HTML
}

type tmplNowPlaying struct {
	Artists string
	Track   string
	State   string
}

// TmplBase …
func TmplBase(w io.Writer, html string) {
	tmpl.ExecuteTemplate(w, "base", tmplBase{
		HTML: template.HTML(html),
	})
}

// TmplPlayer …
func TmplPlayer(w io.Writer, nowPlaying string) {
	tmpl.ExecuteTemplate(w, "player", tmplPlayer{
		NowPlaying: template.HTML(nowPlaying),
	})
}

// TmplNewAuth …
func TmplNewAuth(w io.Writer, signIn string) {
	tmpl.ExecuteTemplate(w, "newAuth", tmplNewAuth{signIn})
}

// TmplNowPlaying …
func TmplNowPlaying(w io.Writer, artists string, track string) {
	tmpl.ExecuteTemplate(w, "nowPlaying", tmplNowPlaying{
		Artists: artists,
		Track:   track,
		State:   "Now Playing",
	})
}

func init() {
	tmpl = template.Must(template.ParseGlob("./templates/*.gohtml"))
}
