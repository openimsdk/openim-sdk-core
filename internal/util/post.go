package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

var (
	BaseURL = ""
	Token   = ""
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
	var reqUrl string
	if host, _ := ctx.Value("apiHost").(string); host != "" {
		reqUrl = host + api
	} else {
		reqUrl = BaseURL + api
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(reqBody))
	if err != nil {
		return errs.ErrInternalServer.Wrap("http.NewRequestWithContext failed " + err.Error())
	}
	request.ContentLength = int64(len(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("operationID", operationID)
	if token, _ := ctx.Value("token").(string); token != "" {
		request.Header.Set("token", token)
	} else {
		request.Header.Set("token", Token)
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
		return errs.ErrInternalServer.Wrap(fmt.Sprintf("api %s json.Unmarshal(`%s`, apiResponse) failed %s", api, string(respBody), err.Error()))
	}
	if baseApi.ErrCode != 0 {
		return errs.NewCodeError(baseApi.ErrCode, baseApi.ErrMsg+" "+baseApi.ErrDlt)
	}
	if resp != nil {
		if err := json.Unmarshal(baseApi.Data, resp); err != nil {
			return errs.ErrInternalServer.Wrap("json.Unmarshal(resp) " + err.Error())
		}
	}
	return nil
}

func CallApi[T any](ctx context.Context, api string, req any) (*T, error) {
	var resp T
	if err := ApiPost(ctx, api, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func GetPageAll[A interface {
	GetPagination() *sdkws.RequestPagination
}, B, C any](ctx context.Context, api string, req A, fn func(resp *B) []C) ([]C, error) {
	if req.GetPagination().ShowNumber == 0 {
		req.GetPagination().ShowNumber = 50
	}
	var res []C
	for i := int32(0); ; i++ {
		req.GetPagination().PageNumber = i
		memberResp, err := CallApi[B](ctx, api, req)
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
