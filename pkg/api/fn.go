package api

import (
	"context"
	"fmt"
	"reflect"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/network"
	"github.com/openimsdk/protocol/sdkws"
)

func newApi[Req, Resp any](api string) Api[Req, Resp] {
	return Api[Req, Resp]{
		api: api,
	}
}

type Api[Req, Resp any] struct {
	api string
}

func (a Api[Req, Resp]) Invoke(ctx context.Context, req *Req) (*Resp, error) {
	var resp Resp
	if err := network.ApiPost(ctx, a.api, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (a Api[Req, Resp]) Execute(ctx context.Context, req *Req) error {
	_, err := a.Invoke(ctx, req)
	return err
}

func (a Api[Req, Resp]) Route() string {
	return a.api
}

// ExtractField is a generic function that extracts a field from the response of a given function.
func ExtractField[A, B, C any](ctx context.Context, fn func(ctx context.Context, req *A) (*B, error), req *A, get func(*B) C) (C, error) {
	resp, err := fn(ctx, req)
	if err != nil {
		var c C
		return c, err
	}
	return get(resp), nil
}

type pagination interface {
	GetPagination() *sdkws.RequestPagination
}

func Page[Req pagination, Resp any, Elem any](ctx context.Context, req Req, api func(ctx context.Context, req Req) (*Resp, error), fn func(*Resp) []Elem) ([]Elem, error) {
	if req.GetPagination() == nil {
		vof := reflect.ValueOf(req)
		for {
			if vof.Kind() == reflect.Ptr {
				vof = vof.Elem()
			} else {
				break
			}
		}
		if vof.Kind() != reflect.Struct {
			return nil, fmt.Errorf("request is not a struct")
		}
		fof := vof.FieldByName("Pagination")
		if !fof.IsValid() {
			return nil, fmt.Errorf("request is not valid Pagination field")
		}
		fof.Set(reflect.ValueOf(&sdkws.RequestPagination{}))
	}
	if req.GetPagination().PageNumber < 0 {
		req.GetPagination().PageNumber = 0
	}
	if req.GetPagination().ShowNumber <= 0 {
		req.GetPagination().ShowNumber = 200
	}
	var result []Elem
	for i := int32(0); ; i++ {
		req.GetPagination().PageNumber = i + 1
		resp, err := api(ctx, req)
		if err != nil {
			return nil, err
		}
		elems := fn(resp)
		result = append(result, elems...)
		if len(elems) < int(req.GetPagination().ShowNumber) {
			break
		}
	}
	return result, nil
}
