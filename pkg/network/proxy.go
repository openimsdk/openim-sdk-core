package network

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/net/proxy"
)

var clientMu sync.RWMutex

func GetHTTPClient() *http.Client { clientMu.RLock(); defer clientMu.RUnlock(); return apiClient }

func SetHTTPConfig(proxyURL string) error {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	dialer := *websocket.DefaultDialer
	dialer.Proxy, dialer.NetDialContext = nil, nil
	if proxyURL != "" {
		u, err := url.Parse(proxyURL)
		if err != nil || u.Host == "" {
			return fmt.Errorf("invalid proxy URL")
		}
		switch strings.ToLower(u.Scheme) {
		case "http", "https":
			transport.Proxy = http.ProxyURL(u)
			dialer.Proxy = http.ProxyURL(u)
		case "socks5", "socks5h":
			var auth *proxy.Auth
			if u.User != nil {
				password, _ := u.User.Password()
				auth = &proxy.Auth{User: u.User.Username(), Password: password}
			}
			d, err := proxy.SOCKS5("tcp", u.Host, auth, proxy.Direct)
			if err != nil {
				return err
			}
			transport.DialContext = func(_ context.Context, network, address string) (net.Conn, error) { return d.Dial(network, address) }
			dialer.NetDialContext = transport.DialContext
		default:
			return fmt.Errorf("unsupported proxy scheme %q", u.Scheme)
		}
	}
	clientMu.Lock()
	apiClient = &http.Client{Timeout: apiClient.Timeout, Transport: transport}
	clientMu.Unlock()
	websocket.DefaultDialer = &dialer
	return nil
}
