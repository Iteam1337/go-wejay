package tmpl

import (
	"html/template"
	"io"
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
)

var (
	tmpl *template.Template
)

type tmplBase struct{ HTML template.HTML }
type tmplNewAuth struct{ SignIn string }
type tmplEmpty struct{ Rooms map[string]int }
type tmplInRoom struct{ Name string }

func Base(w io.Writer, html string) {
	if err := tmpl.ExecuteTemplate(w, "base", tmplBase{template.HTML(html)}); err != nil {
		log.Println(err)
	}
}

func Empty(w io.Writer, available []*message.RefRoom) {
	rooms := make(map[string]int)
	for _, room := range available {
		rooms[room.Id] = int(room.Size)
	}
	if err := tmpl.ExecuteTemplate(w, "empty", tmplEmpty{rooms}); err != nil {
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
