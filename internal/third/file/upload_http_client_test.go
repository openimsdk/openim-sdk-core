//go:build !(js && wasm)

package file

import (
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/network"
)

func TestMultipartUsesConfiguredHTTPClient(t *testing.T) {
	f := newMultipartFixture(t, 8192, 3*8192+17)
	original := network.GetHttpClient()
	if original.Timeout != 10*time.Second {
		t.Fatalf("native HTTP timeout = %s, want 10s", original.Timeout)
	}
	transport := &multipartCountingTransport{RoundTripper: http.DefaultTransport}
	client := &http.Client{Timeout: original.Timeout, Transport: transport}
	network.SetHttpClient(client)
	defer network.SetHttpClient(original)
	for range len(f.sizes) {
		f.release <- struct{}{}
	}
	if _, err := f.uploader.UploadFile(f.ctx, f.request, nil); err != nil {
		t.Fatal(err)
	}
	if got := transport.puts.Load(); got != int32(len(f.sizes)) {
		t.Fatalf("configured client PUT count = %d, want %d", got, len(f.sizes))
	}
	if got := transport.posts.Load(); got < 4 {
		t.Fatalf("configured client API POST count = %d, want at least 4", got)
	}
}

type multipartCountingTransport struct {
	http.RoundTripper
	puts, posts atomic.Int32
}

func (t *multipartCountingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	switch req.Method {
	case http.MethodPut:
		t.puts.Add(1)
	case http.MethodPost:
		t.posts.Add(1)
	}
	return t.RoundTripper.RoundTrip(req)
}
