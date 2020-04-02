package tmpl

import (
	"html/template"
	"io"
	"log"
)

var (
	tmpl *template.Template
)

type tmplBase struct{ HTML template.HTML }
type tmplNewAuth struct{ SignIn string }
type tmplEmpty struct{}
type tmplInRoom struct{ Name string }

func Base(w io.Writer, html string) {
	if err := tmpl.ExecuteTemplate(w, "base", tmplBase{template.HTML(html)}); err != nil {
		log.Println(err)
	}
}

func Empty(w io.Writer) {
	if err := tmpl.ExecuteTemplate(w, "empty", tmplEmpty{}); err != nil {
		log.Println(err)
	}
}

func NewAuth(w io.Writer, signIn string) {
	if err := tmpl.ExecuteTemplate(w, "newAuth", tmplNewAuth{signIn}); err != nil {
		log.Println(err)
	}
}

func InRoom(w io.Writer, room string) {
	if err := tmpl.ExecuteTemplate(w, "room", tmplInRoom{room}); err != nil {
		log.Println(err)
	}
}

func init() {
	tmpl = template.Must(template.ParseGlob("./tmpl/src/*.gohtml"))
}
