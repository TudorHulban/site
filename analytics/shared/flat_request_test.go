package shared

import (
	"encoding/gob"
	"net/http"
)

func ForwardAnalytics(r *http.Request, encoder *gob.Encoder) error {
	flat := FlatRequest{
		RemoteAddr: r.RemoteAddr,
		Host:       r.Host,
		Method:     r.Method,
		URL:        r.URL,
		Header:     r.Header,
	}

	// This streams the raw binary directly across the TCP socket
	return encoder.Encode(flat)
}
