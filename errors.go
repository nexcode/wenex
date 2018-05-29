package wenex

import "errors"

var (
	// ErrHandlerType = Handler must be «http.Handler» or «func(http.ResponseWriter, *http.Request)»
	ErrHandlerType = errors.New("Handler must be «http.Handler» or «func(http.ResponseWriter, *http.Request)»")

	// ErrGetFromConfig = Can't get value from config file (trying to get mismatched types?)
	ErrGetFromConfig = errors.New("Can't get value from config file (trying to get mismatched types?)")

	// ErrConfigListenType = Configuration value «server.http(s).listen» must be a string type
	ErrConfigListenType = errors.New("Configuration value «server.http(s).listen» must be a string type")

	// ErrNoServers = No servers to run. Set «server.http(s).listen» in configuration file
	ErrNoServers = errors.New("No servers to run. Set «server.http(s).listen» in configuration file")
)
