package wenex

import (
	"net/http"
	"net/url"
	"strings"
)

type chain struct {
	handler    []http.Handler
	pattern    []string
	lenPattern int
}

func (h *chain) Chain(handlers ...interface{}) error {
	for _, handler := range handlers {
		switch t := handler.(type) {
		case http.Handler:
			h.handler = append(h.handler, t)
		case func(http.ResponseWriter, *http.Request):
			h.handler = append(h.handler, http.HandlerFunc(t))
		default:
			return ErrHandlerType
		}
	}

	return nil
}

func (c *chain) match(URL *url.URL) bool {
	path := strings.Split(URL.EscapedPath(), "/")
	lenPath := len(path)
	query := URL.Query()

	if c.lenPattern > lenPath {
		return false
	}

	if c.pattern[c.lenPattern-1] != "*" && c.lenPattern < lenPath {
		return false
	}

	for key, value := range c.pattern {
		if len(value) == 0 {
			if len(path[key]) == 0 {
				continue
			}

			return false
		}

		if value[0] == ':' {
			query.Add(value[1:], path[key])
			continue
		}

		if value[0] != '*' && value != path[key] {
			return false
		}
	}

	URL.RawQuery = query.Encode()
	return true
}
