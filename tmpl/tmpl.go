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

type tmplProfile struct {
	WS template.URL
}

// Base …
func Base(w io.Writer, html string) {
	tmpl.ExecuteTemplate(w, "base", tmplBase{template.HTML(html)})
}

// Profile …
func Profile(w io.Writer, ws string) {
	tmpl.ExecuteTemplate(w, "profile", tmplProfile{template.URL(ws)})
}

// NewAuth …
func NewAuth(w io.Writer, signIn string) {
	tmpl.ExecuteTemplate(w, "newAuth", tmplNewAuth{signIn})
}

func init() {
	tmpl = template.Must(template.ParseGlob("./tmpl/src/*.gohtml"))
}
