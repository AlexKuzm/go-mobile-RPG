package websocketsrv

import (
	"context"
	"testing"
	"time"
)

func TestRegisterUser(t *testing.T) {
	s := &WsServer{TokenValidPeriod: time.Hour}
	s.Init()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer s.Shutdown(ctx)
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
}
