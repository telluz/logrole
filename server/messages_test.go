package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/saintpete/logrole/services"
)

func TestInvalidNext(t *testing.T) {
	key := services.NewRandomKey()
	s := &messageListServer{
		SecretKey: key,
	}
	enc, _ := services.Opaque("invalid", key)
	req, _ := http.NewRequest("GET", "/messages?next="+enc, nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Errorf("expected Code to be 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Invalid next page uri") {
		t.Errorf("expected Body to contain error message, got %s", w.Body.String())
	}
}