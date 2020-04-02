package main

import (
	"net/http"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-wejay/cookie"
	"github.com/Iteam1337/go-wejay/tmpl"
	"github.com/Iteam1337/go-wejay/utils"
	"github.com/google/uuid"
)

type routeAuth struct{}

func init() {
	router.auth = routeAuth{}
}

func (route *routeAuth) newAuth(w http.ResponseWriter, r *http.Request) {
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
}

func (route *routeAuth) callback(w http.ResponseWriter, r *http.Request) {
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
}

func (route *routeAuth) signOut(w http.ResponseWriter, r *http.Request) {
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
}
