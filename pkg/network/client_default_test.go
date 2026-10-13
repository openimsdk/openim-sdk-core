//go:build !(js && wasm)

package network

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/tools/errs"
)

func TestApiPostPreservesUploadCancellationCause(t *testing.T) {
	started := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.Copy(io.Discard, r.Body); err != nil {
			t.Error(err)
			return
		}
		close(started)
		<-r.Context().Done()
	}))
	defer server.Close()
	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	ctx = ccontext.WithInfo(ctx, &ccontext.GlobalConfig{IMConfig: &sdk_struct.IMConfig{ApiAddr: server.URL}})
	ctx = ccontext.WithOperationID(ctx, "cancel-signing")
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(context.Canceled)
	done := make(chan error, 1)
	go func() { done <- ApiPost(ctx, "/object/auth_sign", struct{}{}, nil) }()
	<-started
	cancel(sdkerrs.ErrCancelUploadingFile.Wrap())
	err := <-done
	var codeErr errs.CodeError
	if !errors.As(err, &codeErr) || codeErr.Code() != sdkerrs.CancelUploadingFile {
		t.Fatalf("API cancellation cause lost: %v", err)
	}
}
