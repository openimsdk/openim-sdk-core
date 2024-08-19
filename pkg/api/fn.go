package api

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
)

func api[Req, Resp any](api string) func(ctx context.Context, req *Req) (*Resp, error) {
	return func(ctx context.Context, req *Req) (*Resp, error) {
		return util.CallApi[Resp](ctx, api, req)
	}
}
