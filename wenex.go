package wenex

import (
	"context"
	"github.com/nexcode/joneva"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// Wenex struct
type Wenex struct {
	Router  *Router
	Logger  func(string) *log.Logger
	Config  *joneva.Joneva
	servers [2]*http.Server
}

// New return a new Wenex object:
//  configFile: sets default config filename.
//  defaultConfig: contains default configuration parameters.
// Doesn't replace parameters declared in configuration file
// and writes new values to configuration file.
func New(configFile string, defaultConfig map[string]interface{}, logWriter LogWriter) (*Wenex, error) {
	if configFile == "" {
		configFile = "wenex"
	}

	if defaultConfig == nil {
		defaultConfig = DefaultConfig()
	}

	config, err := joneva.New(configFile+".conf", defaultConfig)
	if err != nil {
		return nil, err
	}

	wnx := &Wenex{
		Router: newRouter(),
		Config: config,
	}

	if wnx.Logger, err = NewLogger(config, logWriter); err != nil {
		return nil, err
	}

	if wnx.servers, err = newServer(wnx); err != nil {
		return nil, err
	}

	return wnx, nil
}

// ConnState specifies an optional callback function that is
// called when a client connection changes state. See the
// ConnState type and associated constants for details.
func (wnx *Wenex) ConnState(f func(net.Conn, http.ConnState)) {
	for _, server := range wnx.servers {
		if server != nil {
			server.ConnState = f
		}
	}
}

// Run starts the web server. If an error occurs
// during the operation, the error will be returned.
// This method goes to asleep.
func (wnx *Wenex) Run() error {
	if wnx.servers[0] == nil && wnx.servers[1] == nil {
		return ErrNoServers
	}

	wnx.fixEmptyChain()

	var (
		err [2]error
		wg  sync.WaitGroup
	)

	if wnx.servers[0] != nil {
		wg.Add(1)

		go func() {
			defer wg.Done()
			err[0] = wnx.servers[0].ListenAndServe()
		}()
	}

	if wnx.servers[1] != nil {
		wg.Add(1)

		go func() {
			defer wg.Done()
			err[1] = wnx.servers[1].ListenAndServeTLS("", "")
		}()
	}

	wg.Wait()

	for _, err := range err {
		if err != nil {
			return err
		}
	}

	return nil
}

// Close immediately closes all active net.Listeners and any
// connections in state StateNew, StateActive, or StateIdle. For a
// graceful shutdown, use Shutdown.
//
// Close does not attempt to close (and does not even know about)
// any hijacked connections, such as WebSockets.
//
// Close returns any error returned from closing the Server's
// underlying Listener(s).
func (wnx *Wenex) Close() error {
	var (
		err [2]error
		wg  sync.WaitGroup
	)

	for i := 0; i <= 1; i++ {
		if wnx.servers[i] != nil {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()
				err[i] = wnx.servers[i].Close()
			}(i)
		}
	}

	wg.Wait()

	for _, err := range err {
		if err != nil {
			return err
		}
	}

	return nil
}

// Shutdown gracefully shuts down the server without interrupting any
// active connections. Shutdown works by first closing all open
// listeners, then closing all idle connections, and then waiting
// indefinitely for connections to return to idle and then shut down.
// If the provided context expires before the shutdown is complete,
// Shutdown returns the context's error, otherwise it returns any
// error returned from closing the Server's underlying Listener(s).
//
// When Shutdown is called, Serve, ListenAndServe, and
// ListenAndServeTLS immediately return ErrServerClosed. Make sure the
// program doesn't exit and waits instead for Shutdown to return.
//
// Shutdown does not attempt to close nor wait for hijacked
// connections such as WebSockets. The caller of Shutdown should
// separately notify such long-lived connections of shutdown and wait
// for them to close, if desired. See RegisterOnShutdown for a way to
// register shutdown notification functions.
//
// Once Shutdown has been called on a server, it may not be reused;
// future calls to methods such as Serve will return ErrServerClosed.
func (wnx *Wenex) Shutdown(timeout time.Duration) error {
	var (
		err [2]error
		wg  sync.WaitGroup
	)

	for i := 0; i <= 1; i++ {
		if wnx.servers[i] != nil {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()

				err[i] = wnx.servers[i].Shutdown(ctx)
			}(i)
		}
	}

	wg.Wait()

	for _, err := range err {
		if err != nil {
			return err
		}
	}

	return nil
}

func (wnx *Wenex) fixEmptyChain() {
	for _, method := range wnx.Router.method {
		for _, chain := range method {
			if chain.handler == nil {
				chain.handler = []http.Handler{http.NotFoundHandler()}
			}
		}
	}
}
