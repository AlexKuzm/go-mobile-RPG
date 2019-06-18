package main

import (
	"log"

	"github.com/mishuk-sk/go-pathfinder/server"
)

func main() {
	server := server.NewServer()
	wg := server.ServeWithWebsocket()

	server.WaitShutdown()

	(*wg).Wait()
	log.Printf("DONE!")
}
