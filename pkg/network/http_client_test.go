package network

import (
	"context"

	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
)

func TestName(t *testing.T) {
	var conf ccontext.GlobalConfig
	conf.ApiAddr = "http://127.0.0.1:8080"
	ctx := ccontext.WithInfo(context.Background(), &conf)
	ctx = ccontext.WithOperationID(ctx, "123456")
	var resp any
	if err := ApiPost(ctx, "/test", map[string]any{}, &resp); err != nil {
		t.Log(err)
		return
	}
	t.Log("success")
}
