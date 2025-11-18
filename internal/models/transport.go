package models

import (
	"net/http"
	"examtopics-downloader/internal/constants"
)

// AuthTransport adds Bearer token authentication to requests
type AuthTransport struct {
	Token     string
	Transport http.RoundTripper
}

func (a *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+a.Token)
	return a.Transport.RoundTrip(req)
}

// OptimizedTransport returns a high-performance HTTP transport
func OptimizedTransport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:        constants.MaxIdleConns,
		MaxIdleConnsPerHost: constants.MaxIdleConnsPerHost,
		MaxConnsPerHost:     constants.MaxConnsPerHost,
		IdleConnTimeout:     constants.IdleConnTimeout,
		DisableCompression:  false,
		DisableKeepAlives:   false,
		TLSHandshakeTimeout:   constants.TLSHandshakeTimeout,
		ResponseHeaderTimeout: constants.ResponseHeaderTimeout,
		ExpectContinueTimeout: constants.ExpectContinueTimeout,
	}
}
