package parser

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

type RoutePattern struct {
	Method string
	Host   string
	Path   string
}

var (
	ErrInvalidRoutePattern = errors.New("invalid route pattern")
)

const TotalParts = 2

func ParseRoutePattern(s string) (*RoutePattern, error) {
	rp := &RoutePattern{}

	parts := strings.SplitN(s, " ", TotalParts)

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
		return nil, ErrInvalidRoutePattern
	}

	// Handle host:port
	hostPart := rp.Host
	if h, _, err := net.SplitHostPort(rp.Host); err == nil {
		hostPart = h
	}

	// Strict host check
	if hostPart != "localhost" && !strings.Contains(hostPart, ".") {
		return nil, ErrInvalidRoutePattern
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
