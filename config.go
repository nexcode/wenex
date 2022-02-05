package wenex

import (
	"log"
)

// DefaultConfig returns default configuration options:
//  server.http.listen:   ":http"
//  server.timeout.read:  "30s"
//  server.timeout.write: "30s"
//  server.timeout.idle:  "30s"
//  logger.defaultName:   "wenex"
//  logger.namePrefix:    "log/"
//  logger.usePrefix:     "[!] "
//  logger.useFlag:       log.LstdFlags
func DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"server.http.listen":   ":http",
		"server.timeout.read":  "30s",
		"server.timeout.write": "30s",
		"server.timeout.idle":  "30s",
		"logger.defaultName":   "wenex",
		"logger.namePrefix":    "log/",
		"logger.usePrefix":     "[!] ",
		"logger.useFlag":       log.LstdFlags,
	}
}
