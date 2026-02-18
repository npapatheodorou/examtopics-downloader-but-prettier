package models

import (
	"examtopics-downloader/internal/constants"
	"net/http"
)

// OptimizedTransport returns a high-performance HTTP transport
func OptimizedTransport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:          constants.MaxIdleConns,
		MaxIdleConnsPerHost:   constants.MaxIdleConnsPerHost,
		MaxConnsPerHost:       constants.MaxConnsPerHost,
		IdleConnTimeout:       constants.IdleConnTimeout,
		DisableCompression:    false,
		DisableKeepAlives:     false,
		TLSHandshakeTimeout:   constants.TLSHandshakeTimeout,
		ResponseHeaderTimeout: constants.ResponseHeaderTimeout,
		ExpectContinueTimeout: constants.ExpectContinueTimeout,
	}
}
