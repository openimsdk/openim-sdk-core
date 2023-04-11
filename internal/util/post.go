package util

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"io"
	"net/http"
)

type apiResponse struct {
	ErrCode int             `json:"errCode"`
	ErrMsg  string          `json:"errMsg"`
	ErrDlt  string          `json:"errDlt"`
	Data    json.RawMessage `json:"data"`
}

func ApiPost(ctx context.Context, api string, req, resp any) error {
	operationID, _ := ctx.Value("operationID").(string)
	if operationID == "" {
		return errs.ErrArgs.Wrap("operationID is empty")
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return errs.ErrInternalServer.Wrap("json.Marshal(req) failed " + err.Error())
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, api, bytes.NewReader(reqBody))
	if err != nil {
		return errs.ErrInternalServer.Wrap("http.NewRequestWithContext failed " + err.Error())
	}
	request.ContentLength = int64(len(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("operationID", operationID)
	if token, _ := ctx.Value("token").(string); token != "" {
		request.Header.Set("token", token)
	}
	response, err := new(http.Client).Do(request)
	if err != nil {
		return errs.ErrNetwork.Wrap("ApiPost http.Client.Do failed " + err.Error())
	}
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return errs.ErrNetwork.Wrap("io.ReadAll resp.Body failed " + err.Error())
	}
	var baseApi apiResponse
	if err := json.Unmarshal(respBody, &baseApi); err != nil {
		return errs.ErrInternalServer.Wrap("json.Unmarshal(apiResponse) " + err.Error())
	}
	if resp != nil {
		if err := json.Unmarshal(baseApi.Data, resp); err != nil {
			return errs.ErrInternalServer.Wrap("json.Unmarshal(resp) " + err.Error())
		}
	}
	return nil
}

func CallApi[T any](ctx context.Context, api string, req any) (*T, error) {
	var (
		t    T
		v    any
		resp *T
	)
	if _, ok := any(t).(struct{}); !ok {
		v = &t
		resp = &t
	}
	if err := ApiPost(ctx, api, req, v); err != nil {
		return nil, err
	}
	return resp, nil
}

func GetPageAll[A interface {
	GetPagination() *sdkws.RequestPagination
}, B, C any](ctx context.Context, router string, req A, fn func(resp *B) []C) ([]C, error) {
	if req.GetPagination().ShowNumber == 0 {
		req.GetPagination().ShowNumber = 50
	}
	var res []C
	for i := int32(0); ; i++ {
		req.GetPagination().PageNumber = i
		memberResp, err := CallApi[B](ctx, router, req)
		if err != nil {
			return nil, err
		}
		list := fn(memberResp)
		res = append(res, list...)
		if len(list) < int(req.GetPagination().ShowNumber) {
			break
		}
	}
	return res, nil
}
