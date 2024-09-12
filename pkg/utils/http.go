package utils

import "net/http"

// Get request id from echo context
func GetRequestID(r *http.Request) string {
	return r.Header.Get("X-Request-ID")
}

func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}
