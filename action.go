package main

import "strings"

// Action …
type Action int

// Action enum …
const (
	Unknown Action = iota
	Play
	Pause
	Next
	Previous
	Shuffle
)

func (action Action) String() string {
	return [...]string{
		"unknown",
		"play",
		"pause",
		"next",
		"previous",
		"shuffle",
	}[action]
}

// ActionFromString …
func ActionFromString(action string) Action {
	switch strings.ToLower(action) {
	case "play":
		return Play
	case "pause":
		return Pause
	case "next":
		return Next
	case "previous":
		return Previous
	case "shuffle":
		return Shuffle
	default:
		return Unknown
	}

}
