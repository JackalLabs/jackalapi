package japicore

import (
	"net/http/httptest"
	"testing"
)

func TestHttpGetFileRequest(t *testing.T) {
	w := httptest.NewRecorder()
	gateway := "https://ipfs.io/ipfs/"
	cid := "123"
	_, err := httpGetFileRequest(w, gateway, cid)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
}
