//go:build js

package network

import "net/http"

var httpClient = &http.Client{}

func GetHTTPClient() *http.Client { return httpClient }
func SetHTTPClient(c *http.Client) {
	if c != nil {
		httpClient = c
	}
}
