package wenex

import (
	"log"
	"net/http"
)

// Wenex struct
type Wenex struct {
	Router  *router
	Logger  func(string) *log.Logger
	Config  *Config
	servers [2]*http.Server
}

// New func
func New(configName string, defaultConfig map[string]interface{}) (*Wenex, error) {
	if configName == "" {
		configName = "wenex"
	}

	config, err := newConfig(configName)
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

	if wnx.Logger, err = newLogger(); err != nil {
		return nil, err
	}

	if wnx.servers, err = newServer(wnx); err != nil {
		return nil, err
	}

	return wnx, nil
}

// Run method
func (wnx *Wenex) Run() error {
	if wnx.servers[0] == nil && wnx.servers[1] == nil {
		return ErrNoServers
	}

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
