package api

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
)

func api[Req, Resp any](api string) Api[Req, Resp] {
	return Api[Req, Resp]{
		api: api,
	}
}

type Api[Req, Resp any] struct {
	api string
}

func (a Api[Req, Resp]) Invoke(ctx context.Context, req *Req) (*Resp, error) {
	var resp Resp
	if err := util.ApiPost(ctx, a.api, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (a Api[Req, Resp]) Route() string {
	return a.api
}
