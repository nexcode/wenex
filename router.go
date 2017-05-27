package wenex

import "strings"

func newRouter() *Router {
	return &Router{
		method: make(map[string][]*Chain),
	}
}

type Router struct {
	method map[string][]*Chain
}

func (r *Router) Route(pattern string, methods ...string) *Chain {
	c := &Chain{
		pattern: r.parse(pattern),
	}

	c.lenPattern = len(c.pattern)

	for _, method := range methods {
		r.method[method] = append(r.method[method], c)
	}

	return c
}

func (r *Router) parse(pattern string) []string {
	if len(pattern) == 0 || pattern[0] != '/' {
		pattern = "/" + pattern
	}

	return strings.Split(pattern, "/")
}
