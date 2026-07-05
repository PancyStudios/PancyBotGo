package web

import (
	"net"
	"strings"
)

func isAllowedHost(host string) bool {
	host = stripPort(host)

	if host == "miau.media" || strings.HasSuffix(host, ".miau.media") {
		return true
	}

	switch host {
	case "localhost",
		"127.0.0.1",
		"172.18.0.1",
		"10.66.66.1",
		"10.66.66.4":
		return true
	default:
		return false
	}
}

func stripPort(host string) string {
	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}

	return host
}
