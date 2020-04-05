package main

import (
	"log"
	"net/http"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-wejay/cookie"
	"github.com/Iteam1337/go-wejay/utils"
)

type routeAuth struct{}

func init() {
	router.auth = routeAuth{}
}

func (route *routeAuth) Callback(w http.ResponseWriter, r *http.Request) {
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
		redirect(w, r, routePathEmpty)
	}
}

func (route *routeAuth) SignOut(w http.ResponseWriter, r *http.Request) {
	exists, id, err := exists(r)

	if exists {
		cookie.Expire(w, r)
	}

	if err != nil || !exists {
		redirect(w, r, routePathBase)
		return
	}

	var del message.DeleteUser
	err = updServer.NewRequest(types.IDeleteUser, &message.DeleteUser{UserId: id}, &del)

	if err != nil {
		log.Println(err)
	}

	redirect(w, r, routePathBase)
}
