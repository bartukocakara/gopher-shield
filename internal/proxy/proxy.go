package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/bartukocakara/gopher-shield/internal/resilience"
)

type ShieldProxy struct {
	proxy   *httputil.ReverseProxy
	breaker *resilience.CircuitBreaker
}

func NewShieldProxy(targetURL string, cb *resilience.CircuitBreaker) (*ShieldProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	rp := httputil.NewSingleHostReverseProxy(target)

	return &ShieldProxy{
		proxy:   rp,
		breaker: cb,
	}, nil
}

func (p *ShieldProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var proxyErr error

    p.proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
        proxyErr = e
        // Don't write header here; let the Execute wrapper handle it
    }

    err := p.breaker.Execute(func() error {
        p.proxy.ServeHTTP(w, r)
        return proxyErr
    })

    if err != nil {
        // Only write if the proxy didn't already succeed
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("Shield Active: Service Unavailable"))
    }
}