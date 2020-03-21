package main

import (
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
		id, err := GetIDFromCookie(w, r)
		if err != nil {
			return
		}

		log.Println(id)

		http.NotFound(w, r)

		// // user, ok := users[string(id)]
		// if !ok {
		// 	http.Redirect(w, r, "//"+r.Host+"/new-auth", 307)
		// 	return
		// }

		// w.Header().Set("Content-Type", "text/html")

		// artists, track := user.NowPlaying()
		// var tpl bytes.Buffer
		// TmplNowPlaying(&tpl, artists, track)
		// TmplPlayer(w, tpl.String())
	})

	http.HandleFunc("/action/", func(w http.ResponseWriter, r *http.Request) {
		id, err := GetIDFromCookie(w, r)
		if err != nil {
			json.NewPlayerResponseErr(w, "missing id")
			return
		}

		log.Println("action", id)

		// UserExists <-
		var userExists message.UserExistsResponse
		err = updServer.NewRequest(types.IUserExists, &message.UserExists{UserId: id}, &userExists)
		log.Println("userExists", userExists)
		if !userExists.Ok {
			return
		}

		// Action <-
		action := actionFromString(strings.TrimPrefix(r.URL.Path, "/action/"))
		var actionRes message.ActionResponse
		err = updServer.NewRequest(types.IAction, &message.Action{UserId: id, Action: action}, &actionRes)

		log.Println(actionRes)

		// user, ok := users[string(id)]
		// if !ok {
		// 	json.NewPlayerResponseErr(w, "user not found")
		// 	return
		// }

		// if err := user.RunAction(ActionFromString(action)); err != nil {
		// 	json.NewPlayerResponseErr(w, err.Error())
		// 	return
		// }

		// json.NewPlayerResponse(w, fmt.Sprintf("action %s applied successful", action))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		TmplBase(w, `<a href="/new-auth">new auth</a>`)
	})

	http.HandleFunc("/new-auth", func(w http.ResponseWriter, r *http.Request) {
		var id string
		var cb message.CallbackURLResponse
		var ex message.UserExistsResponse

		id, _ = GetIDFromCookie(w, r)
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
