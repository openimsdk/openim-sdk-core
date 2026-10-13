//go:build js && wasm

package network

import (
	"net/http"
	"sync/atomic"
	"time"
)

var httpClient atomic.Pointer[http.Client]

func init() {
	httpClient.Store(&http.Client{Timeout: time.Second * 30})
}

func GetHttpClient() *http.Client {
	return httpClient.Load()
}

func SetHttpClient(client *http.Client) {

}

func setHTTPTransport(transport *http.Transport) {
	client := *GetHttpClient()
	client.Transport = transport
	httpClient.Store(&client)
}
