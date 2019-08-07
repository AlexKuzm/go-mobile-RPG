package websocketsrv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// ErrUserNotRegistered represents error, when user, who's token doesn't exist in server known tokens,
	// tries to create socket connection
	ErrUserNotRegistered = errors.New("User is not registered to establish websocket connection")
)

// WsServer defines http server, that upgrades connection to websocket and handles it
type WsServer struct {
	Upgrader         websocket.Upgrader
	TokenValidPeriod time.Duration
	server           http.Server
	serveMux         *http.ServeMux
	users            map[string]time.Time
}

// WebsocketHandler represents simple wrapper for function, handling websocket connection
type WebsocketHandler func(*websocket.Conn)

// Init initializes server for future listening
func (s *WsServer) Init() {
	s.server = http.Server{
		Handler:      s.serveMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	s.users = make(map[string]time.Time)
}

// ListenAndServe Starts server listening on provided addr with routes specified by AddWebsocketHandler
func (s *WsServer) ListenAndServe(addr string) error {
	s.server.Addr = addr
	return s.server.ListenAndServe()
}

// RegisterUser adds user token to server-known users for future connection establishment
func (s *WsServer) RegisterUser(token string) error {
	if _, ok := s.users[token]; ok {
		return fmt.Errorf("User with token %s already exists", token)
	}
	s.users[token] = time.Now().Add(s.TokenValidPeriod)
	return nil
}

// AddWebsocketHandler adds new WebsocketHandler to server. Must be added BEFORE server starts serving
func (s *WsServer) AddWebsocketHandler(path string, handleFunc func(*websocket.Conn)) {
	if s.serveMux == nil {
		s.serveMux = http.NewServeMux()
	}
	s.serveMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Unsuccessful connection upgrade: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var token struct {
			Token string `json:"token"`
		}
		err = conn.ReadJSON(&token)
		if err != nil {
			log.Printf("Unable to read JSON FIRST message from ws connection: %v\n", err)
			if err := conn.WriteJSON(struct {
				Err string `json:"err"`
			}{
				"Can't read token from json",
			}); err != nil {
				log.Printf("Unable writing json err to websocket connection: %v", err)
			}
			return
		}
		if err := s.checkUser(token.Token); err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		handleFunc(conn)
	})
}

func (s *WsServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *WsServer) checkUser(token string) error {
	if t, ok := s.users[token]; t.Before(time.Now()) || !ok {
		return ErrUserNotRegistered
	}
	return nil
}
