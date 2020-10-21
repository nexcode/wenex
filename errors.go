package wenex

import "errors"

var (
	// ErrHandlerType = Handler must be «http.Handler» or «func(http.ResponseWriter, *http.Request)»
	ErrHandlerType = errors.New("Handler must be «http.Handler» or «func(http.ResponseWriter, *http.Request)»")

	// ErrConfigValueMismatched = Can't get value from config file (trying to get mismatched types?)
	ErrConfigValueMismatched = errors.New("Can't get value from config file (trying to get mismatched types?)")

	// ErrConfigListenType = Configuration value «server.http(s).listen» must be a string type
	ErrConfigListenType = errors.New("Configuration value «server.http(s).listen» must be a string type")

	// ErrNoServers = No servers to run. Set «server.http(s).listen» in configuration file
	ErrNoServers = errors.New("No servers to run. Set «server.http(s).listen» in configuration file")

	// ErrDefaultLogEmpty = Configuration value «logger.defaultName» must be a non-empty string
	ErrDefaultLogEmpty = errors.New("Configuration value «logger.defaultName» must be a non-empty string")

	// ErrConfigValueNotFound = Failed to get value from config file (not found)
	ErrConfigValueNotFound = errors.New("Failed to get value from config file (not found)")

	// ErrNeedTLSConfigForHTTPS = To create a https lisener, you need to specify the tls config options
	ErrNeedTLSConfigForHTTPS = errors.New("To create a https lisener, you need to specify the tls config options")
)
