package jsonresponses

import (
	"strconv"

	"github.com/Iteam1337/go-protobuf-wejay/message"
)

// ListenResponse …
type ListenResponse struct {
	Ok         bool   `json:"ok,omitempty"`
	Error      string `json:"error,omitempty"`
	Meta       string `json:"meta,omitempty"`
	ActionType string `json:"actionType,omitempty"`
}

// ListenResponseFromMessage …
func ListenResponseFromMessage(m message.ListenResponse) (res ListenResponse) {
	res.Ok = m.Ok
	res.Error = m.Error
	res.ActionType = m.ActionType.String()
	res.Meta = strconv.Itoa(int(m.Meta[0]))
	return
}
