package wenex

import "errors"

var (
	ErrHandlerType      = errors.New("Handler must be «http.Handler» or «func(http.ResponseWriter, *http.Request)»")
	ErrGetFromConfig    = errors.New("Can't get value from config file (trying to get mismatched types?)")
	ErrConfigListenType = errors.New("Configuration value «server.http(s).listen» must be a string type")
	ErrNoServers        = errors.New("No servers to run. Set «server.http(s).listen» in configuration file")
)
