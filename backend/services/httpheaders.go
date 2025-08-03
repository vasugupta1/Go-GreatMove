package services

import (
	"net/http"
)

type HeaderRoundTripper struct {
	Transport http.RoundTripper
	Headers   http.Header
}

func (hrt *HeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	for key, values := range hrt.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return hrt.Transport.RoundTrip(req)
}
