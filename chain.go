package wenex

import (
	"net/http"
)

type Chain struct {
	handler    []http.Handler
	pattern    []string
	strict     bool
	lenPattern int
}

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
