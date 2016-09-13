package wenex

import "errors"

var (
	ErrHandlerType      = errors.New("Handler must be http.Handler or func(http.ResponseWriter, *http.Request)")
	ErrParsePattern     = errors.New("Pattern is not valid")
	ErrGetFromConfig    = errors.New("ErrGetFromConfig")
	ErrConfigListenType = errors.New("\"server.http(s).listen\" must be a string type")
	ErrNoServers        = errors.New("No servers to run. Set \"server.http(s).listen\" in configuration file")
)
