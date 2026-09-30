//go:build js && wasm

package network

import (
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 30,
}

func GetHttpClient() *http.Client {
	return httpClient
}

func SetHttpClient(client *http.Client) {

}
