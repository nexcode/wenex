package wenex

import (
	"testing"
)

func TestNew(t *testing.T) {
	for _, value := range testNewData {
		if _, err := New(value.Name, value.DefaultConfig); err != nil {
			t.Error(err)
		}
	}
}

var testNewData = []struct {
	Name          string
	DefaultConfig map[string]interface{}
}{{
	Name:          "",
	DefaultConfig: nil,
}, {
	Name:          "wenex",
	DefaultConfig: DefaultConfig(),
}, {
	Name: "test/name",
	DefaultConfig: func() map[string]interface{} {
		config := DefaultConfig()
		config["server.http.listen"] = "8080"
		config["server.timeout.idle"] = "10s"
		return config
	}(),
}}
