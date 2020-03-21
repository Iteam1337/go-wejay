package main

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"github.com/Iteam1337/go-wejay/udp"
	"github.com/Iteam1337/go-wejay/utils"
	"github.com/zmb3/spotify"
)

var (
	host        = utils.Getenv("HOST", "localhost")
	port        = utils.Getenv("PORT", "8080")
	addr        = utils.Getenv("ADDR", fmt.Sprintf("%s:%s", host, port))
	udpHost     = utils.Getenv("UDP_SERVER", "localhost:8090")
	wsaddr      = fmt.Sprintf("%s/ws", utils.Getenv("WS_ADDR", fmt.Sprintf("ws://%s", addr)))
	redirectURL = utils.Getenv("REDIRECT_URL", "http://localhost:8080/callback")
	spotifyAuth = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)
	updServer   = udp.Connect(udpHost, true)
)

func main() {
	if strings.HasPrefix(wsaddr, "ws://") && !strings.Contains(wsaddr, "//localhost:") {
		log.Fatalln("websockets will not work unless localhost or wss://")
	}

	ServerListen()
}
