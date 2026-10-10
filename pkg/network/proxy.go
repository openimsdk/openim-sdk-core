package network

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/proxy"
)

func ProxyTransport(proxyURL string) (*http.Transport, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if proxyURL == "" {
		return transport, nil
	}
	u, err := url.Parse(proxyURL)
	if err != nil || u.Host == "" {
		return nil, fmt.Errorf("invalid proxy URL %q", proxyURL)
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		transport.Proxy = http.ProxyURL(u)
	case "socks5", "socks5h":
		dialer, err := proxy.SOCKS5("tcp", u.Host, proxyAuth(u), proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("create SOCKS5 proxy: %w", err)
		}
		transport.Proxy = nil
		transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial(network, address)
		}
	default:
		return nil, fmt.Errorf("unsupported proxy scheme %q", u.Scheme)
	}
	return transport, nil
}

func proxyAuth(u *url.URL) *proxy.Auth {
	if u.User == nil {
		return nil
	}
	p, _ := u.User.Password()
	return &proxy.Auth{User: u.User.Username(), Password: p}
}
