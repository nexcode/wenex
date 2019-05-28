package wenex

import (
	"testing"
)

func TestNew(t *testing.T) {
	if _, err := New("", nil); err != nil {
		t.Error(err)
	}

	if _, err := New("wenex", DefaultConfig()); err != nil {
		t.Error(err)
	}

	config := DefaultConfig()
	config["server.http.listen"] = "8080"
	config["server.timeout.idle"] = "10s"

	if _, err := New("test/name", config); err != nil {
		t.Error(err)
	}
}
