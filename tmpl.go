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

func init() {
	tmpl = template.Must(template.ParseGlob("./templates/*.gohtml"))
}
