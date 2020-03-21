package action

import (
	"strings"

	"github.com/Iteam1337/go-protobuf-wejay/message"
)

// FromString …
func FromString(action string) message.Action_ActionType {
	switch strings.ToLower(action) {
	case "play":
		return message.Action_PLAY
	case "pause":
		return message.Action_PAUSE
	case "next":
		return message.Action_NEXT
	case "previous":
		return message.Action_PREVIOUS
	case "shuffle":
		return message.Action_SHUFFLE
	default:
		return message.Action_UNKNOWN
	}
}
