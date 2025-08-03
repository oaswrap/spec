package parser

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

type RoutePattern struct {
	Method string
	Host   string
	Path   string
}

func ParseRoutePattern(s string) (*RoutePattern, error) {
	rp := &RoutePattern{}

	parts := strings.SplitN(s, " ", 2)

	if len(parts) == 2 && isHTTPMethod(parts[0]) {
		rp.Method = parts[0]
		s = parts[1]
	}
	// If NOT a valid method, keep full input in s!

	if strings.HasPrefix(s, "/") {
		rp.Path = s
		return rp, nil
	}

	slash := strings.Index(s, "/")
	if slash == -1 {
		rp.Host = s
		rp.Path = "/"
	} else {
		rp.Host = s[:slash]
		rp.Path = s[slash:]
	}

	if strings.Contains(rp.Host, " ") {
		return nil, fmt.Errorf("invalid host: contains space")
	}

	// Handle host:port
	hostPart := rp.Host
	if h, _, err := net.SplitHostPort(rp.Host); err == nil {
		hostPart = h
	}

	// Strict host check
	if hostPart != "localhost" && !strings.Contains(hostPart, ".") {
		return nil, fmt.Errorf("invalid host: %q must contain '.' or be 'localhost'", hostPart)
	}

	return rp, nil
}

func isHTTPMethod(s string) bool {
	switch s {
	case http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions,
		http.MethodHead,
		http.MethodTrace,
		http.MethodConnect:
		return true
	default:
		return false
	}
}
