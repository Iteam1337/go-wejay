package tmpl

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

// Base …
func Base(w io.Writer, html string) {
	tmpl.ExecuteTemplate(w, "base", tmplBase{
		HTML: template.HTML(html),
	})
}

// Player …
func Player(w io.Writer, nowPlaying string) {
	tmpl.ExecuteTemplate(w, "player", tmplPlayer{
		NowPlaying: template.HTML(nowPlaying),
	})
}

// NewAuth …
func NewAuth(w io.Writer, signIn string) {
	tmpl.ExecuteTemplate(w, "newAuth", tmplNewAuth{signIn})
}

// NowPlaying …
func NowPlaying(w io.Writer, artists string, track string) {
	tmpl.ExecuteTemplate(w, "nowPlaying", tmplNowPlaying{
		Artists: artists,
		Track:   track,
		State:   "Now Playing",
	})
}

func init() {
	tmpl = template.Must(template.ParseGlob("./src/templates/*.gohtml"))
}
