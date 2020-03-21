package jsonresponses

import (
	"strconv"

	"github.com/Iteam1337/go-protobuf-wejay/message"
)

// ListenResponse …
type ListenResponse struct {
	Ok     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Meta   string `json:"meta,omitempty"`
	Change string `json:"change,omitempty"`
}

// ListenResponseFromMessage …
func ListenResponseFromMessage(m message.ListenResponse) (res ListenResponse) {
	var meta string

	for _, i := range m.Meta {
		meta = meta + strconv.Itoa(int(i))
	}

	res.Ok = m.Ok
	res.Error = m.Error
	res.Change = m.Change.String()
	res.Meta = meta
	return
}
