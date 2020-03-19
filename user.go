package main

import (
	"fmt"
	"log"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// User …
type User struct {
	client *spotify.Client
}

func (u *User) getClient() *spotify.Client {
	return u.client
}

func (u *User) toggleShuffleState() (state bool, err error) {
	client := u.client

	playerState, err := client.PlayerState()

	if err != nil {
		return
	}

	state = !playerState.ShuffleState
	playerState.ShuffleState = state

	return
}

// SetClient …
func (u *User) SetClient(token *oauth2.Token) {
	client := spotifyAuth.NewClient(token)

	u.client = &client

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You are logged in as:", user.ID)

	playerState, err := client.PlayerState()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found your %s (%s)\n", playerState.Device.Type, playerState.Device.Name)
}

// Action …
type Action string

// Action enum …
const (
	Play     Action = "play"
	Pause    Action = "pause"
	Next     Action = "next"
	Previous Action = "previous"
	Shuffle  Action = "shuffle"
)

// RunAction …
func (u *User) RunAction(action Action) (err error) {
	client := u.getClient()
	switch action {
	case Play:
		err = client.Play()
	case Pause:
		err = client.Pause()
	case Next:
		err = client.Next()
	case Previous:
		err = client.Previous()
	case Shuffle:
		state, maybeErr := u.toggleShuffleState()
		if maybeErr != nil {
			err = maybeErr
		} else {
			err = client.Shuffle(state)
		}
	}

	if err != nil {
		log.Print(err)
	}

	return
}
