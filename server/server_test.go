package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	s := NewServer()
	if s.Handler == http.DefaultServeMux {
		t.Logf("WARNING: server uses default ServeMux\n")
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if s.wsSrv.Upgrader.CheckOrigin(req) {
		t.Logf("WARNING: Created websocket server DOESN'T check origin of the request (respods to empty request)\n")
	}
	if s == nil {
		t.Errorf("Failed to create server")
	}
}
func TestServeWithWebsocket(t *testing.T) {
	s := NewServer()
	wg := s.ServeWithWebsocket()
	if wg == nil {
		t.Errorf("Returned WaitGroup is nil. Should be not nil to provide waiting for external routine")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.wsSrv.Shutdown(ctx)
	s.Shutdown(ctx)
}
