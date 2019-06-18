package main

import (
	"log"
	"sync"

	"github.com/mishuk-sk/go-pathfinder/server"
)

func main() {
	server := server.NewServer()
	var wg *sync.WaitGroup
	wg = server.ServeWithWebsocket()

	server.WaitShutdown()

	(*wg).Wait()
	log.Printf("DONE!")
}
