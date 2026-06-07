package analytics

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	// 1. Check standard proxy headers
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can be a comma-separated list;
		// the first one is the client
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// 2. Fallback to RemoteAddr (strip the port number)
	ip, _, errSplit := net.SplitHostPort(r.RemoteAddr)
	if errSplit != nil {
		return r.RemoteAddr
	}

	return ip
}
