package jsonresponses

import (
	"github.com/Iteam1337/go-protobuf-wejay/message"
)

// ListenResponse …
type ListenResponse struct {
	Ok     bool    `json:"ok"`
	Error  string  `json:"error,omitempty"`
	Type   int8    `json:"type,omitempty"`
	Meta   []int16 `json:"meta,omitempty"`
	Change string  `json:"change,omitempty"`
}

// ListenResponseFromMessage …
func ListenResponseFromMessage(m message.ListenResponse) (res ListenResponse) {
	var meta []int16

	for _, i := range m.Meta[1:] {
		meta = append(meta, int16(i))
	}

	res.Ok = m.Ok
	res.Error = m.Error
	res.Change = m.Change.String()
	res.Type = int8(m.Meta[0])
	res.Meta = meta
	return
}
