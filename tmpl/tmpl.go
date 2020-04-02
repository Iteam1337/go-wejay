package tmpl

import (
	"html/template"
	"io"
	"log"
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
}

// Base …
func Base(w io.Writer, html string) {
	if err := tmpl.ExecuteTemplate(w, "base", tmplBase{template.HTML(html)}); err != nil {
		log.Println(err)
	}
}

// Profile …
func Profile(w io.Writer) {
	if err := tmpl.ExecuteTemplate(w, "profile", tmplProfile{}); err != nil {
		log.Println(err)
	}
}

// NewAuth …
func NewAuth(w io.Writer, signIn string) {
	if err := tmpl.ExecuteTemplate(w, "newAuth", tmplNewAuth{signIn}); err != nil {
		log.Println(err)
	}
}

func init() {
	tmpl = template.Must(template.ParseGlob("./tmpl/src/*.gohtml"))
}
