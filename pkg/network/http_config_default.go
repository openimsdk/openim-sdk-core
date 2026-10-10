//go:build !js

package network

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"open_im_sdk/sdk_struct"
)

// SetHTTPConfig configures the shared HTTP client and WebSocket dialer.
func SetHTTPConfig(config sdk_struct.IMConfig) error {
	transport, err := ProxyTransport(config.ProxyURL)
	if err != nil {
		return err
	}
	SetHTTPClient(&http.Client{Timeout: 30 * time.Second, Transport: transport})
	dialer := *websocket.DefaultDialer
	dialer.Proxy = nil
	dialer.NetDialContext = nil
	if config.ProxyURL != "" {
		u, err := url.Parse(config.ProxyURL)
		if err != nil {
			return err
		}
		if strings.HasPrefix(strings.ToLower(u.Scheme), "socks5") {
			dialer.NetDialContext = transport.DialContext
		} else {
			dialer.Proxy = http.ProxyURL(u)
		}
	}
	websocket.DefaultDialer = &dialer
	return nil
}
