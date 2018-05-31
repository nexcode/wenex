package wenex

import (
	"net/http"
	"net/url"
	"strings"
)

func newRouter() *Router {
	return &Router{
		method: make(map[string][]*Chain),
	}
}

type Router struct {
	method map[string][]*Chain
}

func (r *Router) StrictRoute(pattern string, methods ...string) *Chain {
	return r.route(true, pattern, methods)
}

func (r *Router) WeakRoute(pattern string, methods ...string) *Chain {
	return r.route(false, pattern, methods)
}

func (r *Router) route(strict bool, pattern string, methods []string) *Chain {
	c := &Chain{
		strict:  strict,
		pattern: r.parse(pattern),
	}

	c.lenPattern = len(c.pattern)

	for _, method := range methods {
		r.method[method] = append(r.method[method], c)
	}

	return c
}

func (r *Router) parse(pattern string) []string {
	if pattern == "" || pattern[0] != '/' {
		pattern = "/" + pattern
	}

	return strings.Split(pattern, "/")
}

func (r *Router) match(w http.ResponseWriter, re *http.Request) []http.Handler {
	path := strings.Split(re.URL.EscapedPath(), "/")
	lenPath := len(path)

	for _, chain := range r.method[re.Method] {
		if chain.lenPattern > lenPath {
			continue
		}

		if chain.strict && chain.lenPattern < lenPath {
			continue
		}

		query := url.Values{}

		var i int
		var pattern string

		for i, pattern = range chain.pattern {
			if pattern == path[i] {
				continue
			}

			if path[i] != "" {
				if pattern == "*" {
					continue
				}

				if pattern == "" {
					if !chain.strict {
						continue
					}
				} else if pattern[0] == ':' {
					query.Add(pattern[1:], path[i])
					continue
				}
			}

			i--
			break
		}

		if i == chain.lenPattern-1 {
			if len(query) != 0 {
				if re.URL.RawQuery == "" {
					re.URL.RawQuery = query.Encode()
				} else {
					re.URL.RawQuery += "&" + query.Encode()
				}
			}

			return chain.handler
		}
	}

	return []http.Handler{http.NotFoundHandler()}
}
