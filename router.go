package wenex

import "strings"

func newRouter() *router {
	return &router{
		method: make(map[string][]*chain),
	}
}

type router struct {
	method map[string][]*chain
}

func (r *router) Route(pattern string, methods ...string) *chain {
	c := &chain{
		pattern: r.parse(pattern),
	}

	c.lenPattern = len(c.pattern)

	for _, method := range methods {
		r.method[method] = append(r.method[method], c)
	}

	return c
}

func (r *router) parse(pattern string) []string {
	if len(pattern) == 0 || pattern[0] != '/' {
		pattern = "/" + pattern
	}

	return strings.Split(pattern, "/")
}
