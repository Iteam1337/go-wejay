package main

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"

	"github.com/Iteam1337/go-wejay/udp"
	"github.com/Iteam1337/go-wejay/utils"
)

var (
	host      = utils.Getenv("HOST", "localhost")
	port      = utils.Getenv("PORT", "8080")
	addr      = utils.Getenv("ADDR", fmt.Sprintf("%s:%s", host, port))
	udpHost   = utils.Getenv("UDP_SERVER", "localhost:8090")
	updServer = udp.Connect(udpHost)
)

func main() {
	serverListen()
}
