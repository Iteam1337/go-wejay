package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-wejay/cookie"
	"github.com/Iteam1337/go-wejay/tmpl"
	"github.com/Iteam1337/go-wejay/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// ServerListen …
func ServerListen() {
	router := mux.NewRouter()

	router.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		var res message.NewUserResponse

		id := cookie.CreateAndSet(w, r)
		code, _ := utils.ParseRequest(id, r)

		if err := updServer.NewRequest(
			types.INewUser,
			&message.NewUser{UserId: id, Code: code},
			&res,
		); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Redirect(w, r, "//"+r.Host+"/profile", 307)
		}
	})

	router.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "//"+r.Host+"/", 307)
	})

	router.HandleFunc("/room/{room}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.Header().Set("Content-Type", "application/json")
		if _, e := w.Write([]byte(`{"path":"/room/` + vars["room"] + `"}`)); e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}
	})

	router.HandleFunc("/room/leave", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, e := w.Write([]byte(`{"path":"/room/leave"}`)); e != nil {
			http.Error(w, e.Error(), http.StatusInternalServerError)
		}
	})

	router.HandleFunc("/sign-out", func(w http.ResponseWriter, r *http.Request) {
		exists, id, err := exists(r)

		if exists {
			cookie.Expire(w, r)
		}

		if err != nil || !exists {
			http.Redirect(w, r, "//"+r.Host+"/", 307)
			return
		}

		var del message.DeleteUser
		err = updServer.NewRequest(types.IDeleteUser, &message.DeleteUser{UserId: id}, &del)

		if err != nil {
			http.Redirect(w, r, "//"+r.Host+"/", 307)
			return
		}

		http.Redirect(w, r, "//"+r.Host+"/", 307)
	})

	router.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
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
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `<a href="/new-auth">new auth</a>`
		exists, _, err := exists(r)

		if err != nil || !exists {
			tmpl.Base(w, html)
			return
		}

		http.Redirect(w, r, "//"+r.Host+"/profile", 307)
	})

	router.HandleFunc("/new-auth", func(w http.ResponseWriter, r *http.Request) {
		exists, _, err := exists(r)

		if err == nil && exists {
			http.Redirect(w, r, "//"+r.Host+"/profile", 307)
			return
		}

		var cb message.CallbackURLResponse
		if err := updServer.NewRequest(
			types.ICallbackURL,
			&message.CallbackURL{UserId: uuid.New().String()},
			&cb,
		); err != nil {
			http.Error(w, "Couldn't get callback-url", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		tmpl.NewAuth(w, cb.Url)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", router)

	log.Printf("Listen on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
