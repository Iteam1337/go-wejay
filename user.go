package main

import (
	"fmt"
	"log"
	"strings"

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
}

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

// NowPlaying …
func (u User) NowPlaying() (html string) {
	client := u.getClient()

	currentlyPlaying, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	item := currentlyPlaying.Item
	if item == nil {
		return
	}

	var artists []string
	for _, key := range item.Artists {
		artists = append(artists, key.Name)
	}
	artist := strings.Join(artists, ", ")

	html = fmt.Sprintf(`
	<div>
		<strong>Now playing</strong>
		<p>
			<span class="track">%s</span> ­- <span class="artists">%s</span>
		</p>
	</div>
	`, item.SimpleTrack.Name, artist)

	return
}
