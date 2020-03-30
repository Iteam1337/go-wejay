package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-wejay/action"
	"github.com/Iteam1337/go-wejay/cookie"
	"github.com/Iteam1337/go-wejay/jsonresponses"
	"github.com/Iteam1337/go-wejay/tmpl"
	"github.com/Iteam1337/go-wejay/utils"
	"github.com/google/uuid"
)

// ServerListen …
func ServerListen() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/ws", wsListen)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
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
			http.Redirect(w, r, "//"+r.Host+"/player", 307)
		}
	})

	http.HandleFunc("/remove-user", func(w http.ResponseWriter, r *http.Request) {
		id, err := cookie.GetID(r)
		if err != nil {
			http.Redirect(w, r, "//"+r.Host+"/", 307)
			return
		}

		cookie.Expire(w, r)

		var del message.DeleteUser
		err = updServer.NewRequest(types.IDeleteUser, &message.DeleteUser{UserId: id}, &del)

		if err != nil {
			http.Redirect(w, r, "//"+r.Host+"/", 307)
			return
		}

		http.Redirect(w, r, "//"+r.Host+"/", 307)
	})

	http.HandleFunc("/player/", func(w http.ResponseWriter, r *http.Request) {
		id, err := cookie.GetID(r)
		if err != nil {
			http.Redirect(w, r, "//"+r.Host+"/", 307)
			return
		}

		// UserExists <-
		var userExists message.UserExistsResponse
		if err = updServer.NewRequest(types.IUserExists, &message.UserExists{UserId: id}, &userExists); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !userExists.Ok || !userExists.Exists {
			http.Redirect(w, r, "//"+r.Host+"/new-auth", 307)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		// NowPlaying <-
		var nowPlaying message.NowPlayingResponse
		if err := updServer.NewRequest(types.INowPlaying, &message.NowPlaying{UserId: id}, &nowPlaying); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		var artist string = ""
		var trackName string = ""

		track := nowPlaying.Track
		if track != nil {
			var artists []string
			for _, key := range track.Artists {
				artists = append(artists, key.Name)
			}
			artist = strings.Join(artists, ", ")
			trackName = track.Name
		}

		var tpl bytes.Buffer

		tmpl.NowPlaying(&tpl, artist, trackName)
		tmpl.Player(w, tpl.String(), fmt.Sprintf("ws://%s/ws", addr))
	})

	http.HandleFunc("/action/", func(w http.ResponseWriter, r *http.Request) {
		id, err := cookie.GetID(r)
		if err != nil {
			jsonresponses.NewPlayerResponseErr(w, "missing id")
			return
		}

		// UserExists <-
		var userExists message.UserExistsResponse
		updServer.NewRequest(types.IUserExists, &message.UserExists{UserId: id}, &userExists)
		if !userExists.Ok || !userExists.Exists {
			jsonresponses.NewPlayerResponseErr(w, "couldn't find user")
			return
		}

		// Action <-
		actionStr := strings.TrimPrefix(r.URL.Path, "/action/")
		action := action.FromString(actionStr)
		var actionRes message.ActionResponse
		updServer.NewRequest(types.IAction, &message.Action{UserId: id, Action: action}, &actionRes)
		if !actionRes.Ok {
			jsonresponses.NewPlayerResponseErr(w, actionRes.Error)
			return
		}

		jsonresponses.NewPlayerResponse(w, fmt.Sprintf("action %s applied successful", actionStr))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Base(w, `<a href="/new-auth">new auth</a>`)
	})

	http.HandleFunc("/new-auth", func(w http.ResponseWriter, r *http.Request) {
		var id string
		var cb message.CallbackURLResponse
		var ex message.UserExistsResponse

		id, _ = cookie.GetID(r)
		if id != "" {
			updServer.NewRequest(
				types.IUserExists,
				&message.UserExists{UserId: id},
				&ex,
			)

			if ex.Exists {
				http.Redirect(w, r, "//"+r.Host+"/player", 307)
				return
			}
		}

		id = uuid.New().String()
		if err := updServer.NewRequest(
			types.ICallbackURL,
			&message.CallbackURL{UserId: id},
			&cb,
		); err != nil {
			http.Error(w, "Couldn't get callback-url", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		tmpl.NewAuth(w, cb.Url)
	})

	log.Printf("Listen on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
