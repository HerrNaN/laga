package web

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ServeProxy(rawWebServerAddress string) (http.HandlerFunc, error) {
	target, err := url.Parse(rawWebServerAddress)
	if err != nil {
		return nil, fmt.Errorf("parse web server address: %w", err)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(target)
			pr.In.Host = target.Host
		},
	}

	return proxy.ServeHTTP, nil
}
