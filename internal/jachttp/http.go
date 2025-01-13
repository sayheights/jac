package jachttp

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

var Client = &http.Client{Timeout: time.Minute * 30, Transport: transport}

var transport = newTransport()

func newTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func NewClientWithTLS(tlsConf *tls.Config) *http.Client {
	t := newTransport()
	t.TLSClientConfig = tlsConf
	t.MaxIdleConns = 20

	return &http.Client{Timeout: time.Second * 30, Transport: t}
}

func SetHeaders(req *http.Request, headers http.Header) {
	for k, vals := range headers {
		for _, val := range vals {
			req.Header.Add(k, val)
		}
	}
}
