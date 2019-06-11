package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ClientManager struct {
	clients     map[uuid.UUID]bool
	subscribe   chan *Client
	unsubscribe chan *Client
}

type Client struct {
	id     uuid.UUID
	socket *websocket.Conn
	send   chan WsEvent
}

func wsNewClient(w http.ResponseWriter, r *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
		return
	}

}
func main() {
	router := mux.NewRouter()
	wsRouter := router.PathPrefix("/ws").Subrouter()
	wsRouter.HandleFunc("", wsNewClient)
	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				json.NewEncoder(w).Encode("Smth wrong")
			}
			for {
				msgType, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println(err)
				}
				var message WsEvent
				if msgType == websocket.TextMessage {
					err := json.Unmarshal(msg, &message)
					if err != nil {
						log.Println(err)
					}
				}
				// Print the message to the console
				fmt.Printf("%s sent: %v\n", conn.RemoteAddr(), message)

			}
		})
		http.ListenAndServe(":8080", nil)*/
}
