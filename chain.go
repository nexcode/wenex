package wenex

import (
	"net/http"
)

// Chain struct
type Chain struct {
	handler    []http.Handler
	pattern    []string
	strict     bool
	lenPattern int
}

// Chain adds a http.Handler or func(http.ResponseWriter, *http.Request) to the chain.
// It will be called on http request when the current router is selected.
func (c *Chain) Chain(handlers ...interface{}) error {
	for _, handler := range handlers {
		switch t := handler.(type) {
		case http.Handler:
			c.handler = append(c.handler, t)
		case func(http.ResponseWriter, *http.Request):
			c.handler = append(c.handler, http.HandlerFunc(t))
		default:
			return ErrHandlerType
		}
	}

	return nil
}

// MustChain same as Chain, but causes panic
func (c *Chain) MustChain(handlers ...interface{}) {
	if err := c.Chain(handlers...); err != nil {
		panic(err)
	}
}
