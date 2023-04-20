package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"io"
	"net/http"
	"time"
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

func ApiPost(ctx context.Context, api string, req, resp any) (err error) {
	operationID, _ := ctx.Value("operationID").(string)
	if operationID == "" {
		err := errs.ErrArgs.Wrap("operationID is empty")
		log.ZError(ctx, "ApiRequest", err, "type", "ctx not set operationID")
		return err
	}
	defer func(start time.Time) {
		end := time.Now()
		if err == nil {
			log.ZDebug(ctx, "CallApi", "api", api, "use", "state", "success", time.Duration(end.UnixNano()-start.UnixNano()))
		} else {
			log.ZError(ctx, "CallApi", err, "api", api, "use", "state", "failed", time.Duration(end.UnixNano()-start.UnixNano()))
		}
	}(time.Now())
	reqBody, err := json.Marshal(req)
	if err != nil {
		log.ZError(ctx, "ApiRequest", err, "type", "json.Marshal(req) failed")
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
		log.ZError(ctx, "ApiRequest", err, "type", "http.NewRequestWithContext failed")
		return errs.ErrInternalServer.Wrap("http.NewRequestWithContext failed " + err.Error())
	}
	log.ZDebug(ctx, "ApiRequest", "url", reqUrl, "body", string(reqBody))
	request.ContentLength = int64(len(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("operationID", operationID)
	if token, _ := ctx.Value("token").(string); token != "" {
		request.Header.Set("token", token)
		log.ZDebug(ctx, "ApiRequestToken", "source", "context", "token", token)
	} else {
		request.Header.Set("token", Token)
		log.ZDebug(ctx, "ApiRequestToken", "source", "global", "token", token)
	}
	response, err := new(http.Client).Do(request)
	if err != nil {
		log.ZError(ctx, "ApiRequest", err, "type", "network error")
		return errs.ErrNetwork.Wrap("ApiPost http.Client.Do failed " + err.Error())
	}
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.ZError(ctx, "ApiResponse", err, "type", "read body", "status", response.Status)
		return errs.ErrNetwork.Wrap("io.ReadAll resp.Body failed " + err.Error())
	}
	log.ZDebug(ctx, "ApiResponse", "url", reqUrl, "status", response.Status, "body", string(respBody))
	var baseApi apiResponse
	if err := json.Unmarshal(respBody, &baseApi); err != nil {
		log.ZError(ctx, "ApiResponse", err, "type", "api code parse")
		return errs.ErrInternalServer.Wrap(fmt.Sprintf("api %s json.Unmarshal(`%s`, apiResponse) failed %s", api, string(respBody), err.Error()))
	}
	if baseApi.ErrCode != 0 {
		err := errs.NewCodeError(baseApi.ErrCode, baseApi.ErrMsg+" "+baseApi.ErrDlt)
		log.ZError(ctx, "ApiResponse", err, "type", "api code error", "msg", baseApi.ErrMsg, "dlt", baseApi.ErrDlt)
		return err
	}
	if resp != nil {
		if err := json.Unmarshal(baseApi.Data, resp); err != nil {
			log.ZError(ctx, "ApiResponse", err, "type", "api data parse", "data", string(baseApi.Data), "bind", fmt.Sprintf("%T", resp))
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
