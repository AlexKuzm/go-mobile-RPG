package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mishuk-sk/go-pathfinder/server/websocketsrv"
)

func TestGetToken(t *testing.T) {
	wsSrv := &websocketsrv.WsServer{}
	go func() {
		wsSrv.Init()
		wsSrv.ListenAndServe(":9876")
	}()

	fn := getToken(wsSrv)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	fn(rec, req)
	respBody := make([]byte, 100)
	if _, err := rec.Result().Body.Read(respBody); err != nil && err != io.EOF {
		t.Skipf("Skipping test. Can't read response body: %v\n", err)
	}
	if stat := rec.Result().StatusCode; stat != http.StatusOK {
		t.Errorf("Error returning response. Status code is: %v; expected 200.\n Response body: %s\n", stat, respBody)
	} else if len(respBody) == 0 {
		t.Errorf("Error: Empty body returned")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := wsSrv.Shutdown(ctx);  err != nil {
		t.Errorf("Error shutting down the server (handlers): %v\n", err)
	}
}
