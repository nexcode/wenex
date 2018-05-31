package wenex

import (
	"log"
	"net"
	"net/http"
)

// Wenex struct
type Wenex struct {
	Router  *Router
	Logger  func(string) *log.Logger
	Config  *Config
	servers [2]*http.Server
}

// New return a new Wenex object:
//  defaultName: sets default config filename and default log filename.
//  defaultConfig: contains default configuration parameters.
// Doesn't replace parameters declared in configuration file
// and writes new values to configuration file.
func New(defaultName string, defaultConfig map[string]interface{}) (*Wenex, error) {
	if defaultName == "" {
		defaultName = "wenex"
	}

	config, err := newConfig(defaultName)
	if err != nil {
		return nil, err
	}

	if defaultConfig != nil {
		var needSave bool

		for key, value := range defaultConfig {
			if config.Get(key) == nil {
				config.Set(key, value)
				needSave = true
			}
		}

		if needSave {
			if err = config.Save(); err != nil {
				return nil, err
			}
		}
	}

	wnx := &Wenex{
		Router: newRouter(),
		Config: config,
	}

	if wnx.Logger, err = newLogger(wnx, defaultName); err != nil {
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

	wnx.chainValidation()
	stop := make(chan error)

	if wnx.servers[0] != nil {
		go func() {
			stop <- wnx.servers[0].ListenAndServe()
		}()
	}

	if wnx.servers[1] != nil {
		go func() {
			crt, err := wnx.Config.String("server.https.crt")
			if err != nil {
				stop <- err
			}

			key, err := wnx.Config.String("server.https.key")
			if err != nil {
				stop <- err
			}

			stop <- wnx.servers[1].ListenAndServeTLS(crt, key)
		}()
	}

	return <-stop
}

func (wnx *Wenex) chainValidation() {
	for _, method := range wnx.Router.method {
		for _, chain := range method {
			if chain.handler == nil {
				chain.handler = []http.Handler{http.NotFoundHandler()}
			}
		}
	}
}
