package wenex

import (
	"net/http"
	"time"
)

func newServer(wnx *Wenex) ([2]*http.Server, error) {
	var servers [2]*http.Server

	tmp, err := wnx.Config.String("server.timeout.read")
	if err != nil {
		return servers, err
	}

	rTimeout, err := time.ParseDuration(tmp)
	if err != nil {
		return servers, err
	}

	tmp, err = wnx.Config.String("server.timeout.write")
	if err != nil {
		return servers, err
	}

	wTimeout, err := time.ParseDuration(tmp)
	if err != nil {
		return servers, err
	}

	tmp, err = wnx.Config.String("server.timeout.idle")
	if err != nil {
		return servers, err
	}

	idleTimeout, err := time.ParseDuration(tmp)
	if err != nil {
		return servers, err
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, value := range wnx.Router.method[r.Method] {
			if value.match(r.URL) {
				run := newRun(w, r, value.handler)

				for {
					if !run.Next() {
						return
					}
				}
			}
		}

		http.NotFound(w, r)
	})

	if addr := wnx.Config.Get("server.http.listen"); addr != nil {
		addr, ok := addr.(string)
		if !ok {
			return servers, ErrConfigListenType
		}

		servers[0] = &http.Server{
			Addr:         addr,
			ErrorLog:     wnx.Logger(""),
			Handler:      handler,
			ReadTimeout:  rTimeout,
			WriteTimeout: wTimeout,
			IdleTimeout:  idleTimeout,
		}
	}

	if addr := wnx.Config.Get("server.https.listen"); addr != nil {
		addr, ok := addr.(string)
		if !ok {
			return servers, ErrConfigListenType
		}

		servers[1] = &http.Server{
			Addr:         addr,
			ErrorLog:     wnx.Logger(""),
			Handler:      handler,
			ReadTimeout:  rTimeout,
			WriteTimeout: wTimeout,
			IdleTimeout:  idleTimeout,
		}
	}

	return servers, nil
}
