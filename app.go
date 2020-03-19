package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/Iteam1337/go-wejay/utils"
	"github.com/zmb3/spotify"
)

var (
	addr        = utils.Getenv("ADDR", ":8080")
	redirectURL = utils.Getenv("REDIRECT_URL", "http://localhost:8080/callback")
	spotifyAuth = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	users       = make(map[string]*User)
)

func main() {
	ServerListen()
}
