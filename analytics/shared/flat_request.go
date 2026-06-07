package shared

import "net/url"

type FlatRequest struct {
	RemoteAddr string
	Host       string
	Method     string
	URL        *url.URL
	Header     map[string][]string
}
