package wenex

import (
	"context"
	"net/http"
)

func newRun(w http.ResponseWriter, r *http.Request, handler []http.Handler) *Run {
	run := &Run{
		rWriter: w,
		handler: handler,
	}

	run.request = r.WithContext(context.WithValue(r.Context(), ctxRun, run))
	return run
}

// Run struct
type Run struct {
	rWriter http.ResponseWriter
	request *http.Request
	handler []http.Handler
}

// Next method
func (r *Run) Next() bool {
	if len(r.handler) == 0 {
		return false
	}

	handler := r.handler[0]
	r.handler = r.handler[1:]
	handler.ServeHTTP(r.rWriter, r.request)
	return true
}

// Break method
func (r *Run) Break() {
	r.handler = r.handler[:0]
}
