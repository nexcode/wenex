package wenex

import (
	"net/http"
	"testing"
)

func TestChain(t *testing.T) {
	wnx, err := New("", nil)
	if err != nil {
		t.Error(err)
	}

	if err := wnx.Router.StrictRoute("/", "GET").Chain(nil); err != ErrHandlerType {
		t.Error("Incorrect handler type validation")
	}

	if err := wnx.Router.StrictRoute("/", "GET").Chain(""); err != ErrHandlerType {
		t.Error("Incorrect handler type validation")
	}

	if err := wnx.Router.StrictRoute("/", "GET").Chain(123); err != ErrHandlerType {
		t.Error("Incorrect handler type validation")
	}

	f := func(w http.ResponseWriter, r *http.Request) {}
	if err := wnx.Router.StrictRoute("/", "GET").Chain(f); err != nil {
		t.Error(err)
	}

	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	if err := wnx.Router.StrictRoute("/", "GET").Chain(hf); err != nil {
		t.Error(err)
	}
}
