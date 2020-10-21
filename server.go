package wenex

import (
	"compress/gzip"
	"crypto/tls"
	"net/http"
	"strings"
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
		run := newRun(w, r, wnx.Router.match(r))

		for {
			if !run.Next() {
				break
			}
		}
	})

	gzipEnable, err := wnx.Config.Bool("server.gzip.enable")
	if err != nil && err != ErrConfigValueNotFound {
		return servers, err
	}

	if gzipEnable {
		gzipLevel, err := wnx.Config.Float64("server.gzip.level")
		if err != nil {
			return servers, err
		}

		h := handler

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.Header().Set("Content-Encoding", "gzip")

				gz, err := gzip.NewWriterLevel(w, int(gzipLevel))
				if err != nil {
					panic(err)
				}

				defer gz.Close()
				w = gzipResponseWriter{ResponseWriter: w, Writer: gz}
			}

			h.ServeHTTP(w, r)
		})
	}

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

		if servers[1].TLSConfig, err = stringCert(wnx, "server.https.stringCert"); err != nil {
			return servers, err
		}

		if servers[1].TLSConfig != nil {
			return servers, nil
		}

		if servers[1].TLSConfig, err = loadCert(wnx, "server.https.loadCert"); err != nil {
			return servers, err
		}

		if servers[1].TLSConfig != nil {
			return servers, nil
		}

		certManager, err := autoCert(wnx, "server.https.autoCert")
		if err != nil {
			return servers, err
		}

		if certManager == nil {
			return servers, ErrNeedTLSConfigForHTTPS
		}

		servers[0].Handler = certManager.HTTPHandler(handler)

		servers[1].TLSConfig = &tls.Config{
			GetCertificate: certManager.GetCertificate,
		}
	}

	return servers, nil
}
