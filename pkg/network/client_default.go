//go:build !js

package network

import (
	"net/http"
	"sync/atomic"
	"time"
)

var httpClient atomic.Pointer[http.Client]

func init()                        { httpClient.Store(&http.Client{Timeout: 30 * time.Second}) }
func GetHTTPClient() *http.Client  { return httpClient.Load() }
func SetHTTPClient(c *http.Client) { httpClient.Store(c) }
