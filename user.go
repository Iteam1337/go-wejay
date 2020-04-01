package main

import (
	"net/http"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-wejay/cookie"
)

func exists(r *http.Request) (exists bool, id string, err error) {
	var ex message.UserExistsResponse

	id, _ = cookie.GetID(r)

	if id != "" {
		err = updServer.NewRequest(
			types.IUserExists,
			&message.UserExists{UserId: id},
			&ex,
		)

		exists = ex.Ok && ex.Exists
	}

	return
}
