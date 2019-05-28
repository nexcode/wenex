package wenex

import (
	"net/http"
	"testing"
)

func TestNext(t *testing.T) {
	handlers := []http.Handler{
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	}

	run := &Run{
		handler: handlers,
	}

	var count int
	for run.Next() {
		count++
	}

	if count != 3 {
		t.Error("Not all handlers are processed")
	}
}
