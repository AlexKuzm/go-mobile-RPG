package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/mishuk-sk/go-pathfinder/server/websocketsrv"

	"github.com/gorilla/websocket"
)

// WithWS defines both regular http server for REST API
// and websocketsrv.WsServer for websockets
type WithWS struct {
	http.Server
	wsSrv websocketsrv.WsServer
}

// NewServer creates new server instance and initializes routes for it
func NewServer() *WithWS {
	s := &WithWS{
		http.Server{
			// TODO add port escaping
			Addr:         ":8080",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		websocketsrv.WsServer{
			Upgrader: websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				// TODO update CheckOrigin to validate only OUR client headers
				CheckOrigin: func(r *http.Request) bool { return true },
			},
			TokenValidPeriod: time.Hour,
		},
	}
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("hello")) })
	s.wsSrv.AddWebsocketHandler("/", func(conn *websocket.Conn) { conn.WriteMessage(websocket.TextMessage, []byte("world")) })
	s.Handler = router

	return s
}

// WaitShutdown is used to wait until OS interrupt signal to gracefully shutdown servers
func (s *WithWS) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)
	sig := <-irqSig
	log.Printf("Shutdown request by system signal: %v\n", sig)
	// TODO ... I've just Copy Pasted this part... WTF is context???
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.wsSrv.Shutdown(ctx); err != nil {
		log.Printf("Websocket shutdown request err: %v\n", err)
	}
	if err := s.Shutdown(ctx); err != nil {
		log.Printf("Shutdown request err: %v\n", err)
	}
}

// ServeWithWebsocket starts both servers (regular http and websocket one)
// it returns *sync.Waitgroup, that purpose is to define when both server starts
// finish their work
func (s *WithWS) ServeWithWebsocket() *sync.WaitGroup {
	// TODO Add port escaping
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		if err := s.wsSrv.ListenAndServe(":8090"); err != nil {
			log.Println("Enable to start websocket server: " + err.Error())
		}
		wg.Done()
	}()
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Println("Enable to start http server: " + err.Error())
		}
		wg.Done()
	}()
	return &wg
}
