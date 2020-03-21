package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-protobuf-wejay/types"
	json "github.com/Iteam1337/go-wejay/jsonResponses"
	"github.com/Iteam1337/go-wejay/utils"
	"github.com/google/uuid"
)

// ServerListen …
func ServerListen() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		var res message.NewUserResponse

		id := CreateAndSetCookie(w, r)
		code, err := utils.ParseRequest(id, r)

		if err = updServer.NewRequest(
			types.INewUser,
			&message.NewUser{UserId: id, Code: code},
			&res,
		); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Redirect(w, r, "//"+r.Host+"/player", 307)
		}
		log.Println(res)
	})

	http.HandleFunc("/player/", func(w http.ResponseWriter, r *http.Request) {
		id, err := GetIDFromCookie(r)
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

		log.Println("player", "userExists", id, userExists.Ok, userExists.Exists)
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

		log.Println("player", "nowPlaying", id, nowPlaying.Track)

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

		TmplNowPlaying(&tpl, artist, trackName)
		TmplPlayer(w, tpl.String())
	})

	http.HandleFunc("/action/", func(w http.ResponseWriter, r *http.Request) {
		id, err := GetIDFromCookie(r)
		if err != nil {
			json.NewPlayerResponseErr(w, "missing id")
			return
		}

		// UserExists <-
		var userExists message.UserExistsResponse
		err = updServer.NewRequest(types.IUserExists, &message.UserExists{UserId: id}, &userExists)
		log.Println("action", "userExists", id, userExists.Ok, userExists.Exists)
		if !userExists.Ok || !userExists.Exists {
			json.NewPlayerResponseErr(w, "couldn't find user")
			return
		}

		// Action <-
		actionStr := strings.TrimPrefix(r.URL.Path, "/action/")
		action := actionFromString(actionStr)
		var actionRes message.ActionResponse
		err = updServer.NewRequest(types.IAction, &message.Action{UserId: id, Action: action}, &actionRes)
		if !actionRes.Ok {
			json.NewPlayerResponseErr(w, actionRes.Error)
			return
		}

		json.NewPlayerResponse(w, fmt.Sprintf("action %s applied successful", actionStr))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		TmplBase(w, `<a href="/new-auth">new auth</a>`)
	})

	http.HandleFunc("/new-auth", func(w http.ResponseWriter, r *http.Request) {
		var id string
		var cb message.CallbackURLResponse
		var ex message.UserExistsResponse

		id, _ = GetIDFromCookie(r)
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
		TmplNewAuth(w, cb.Url)
	})

	log.Printf("Listen on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
