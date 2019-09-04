package server

import (
	"log"
	"net/http"

	"github.com/mishuk-sk/go-pathfinder/server/websocketsrv"

	"github.com/google/uuid"
)

func getToken(wsSrv *websocketsrv.WsServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.NewUUID()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := wsSrv.RegisterUser(id.String()); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(id.String()))
	}
}
