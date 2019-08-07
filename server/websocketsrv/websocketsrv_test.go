package websocketsrv

import (
	"context"
	"testing"
	"time"
)

func TestRegisterUser(t *testing.T) {
	s := &WsServer{TokenValidPeriod: time.Hour}
	s.Init()
	token := "token"
	if err := s.RegisterUser(token); err != nil {
		t.Errorf("Error registering user: %v\n", err)
	}
	if err := s.RegisterUser(token); err == nil {
		t.Errorf("Somehow 2 users with same token were registered\n")
	}
	if err := s.checkUser(token); err != nil {
		t.Errorf("Can't fin registered user with token %v\n", token)
	}
	if err := s.checkUser(token + "smth"); err == nil {
		t.Errorf("Thought, that user with token %v exists, but it doesn't", token+"smth")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err :=  s.Shutdown(ctx); err != nil {
		t.Errorf("Error shutting down the server: %v\n", err)
	}
}
