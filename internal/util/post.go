// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/sdkerrs"
	"time"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
)

//var (
//	BaseURL = ""
//	Token   = ""
//)

type apiResponse struct {
	ErrCode int             `json:"errCode"`
	ErrMsg  string          `json:"errMsg"`
	ErrDlt  string          `json:"errDlt"`
	Data    json.RawMessage `json:"data"`
}

func ApiPost(ctx context.Context, api string, req, resp any) (err error) {
	operationID, _ := ctx.Value("operationID").(string)
	if operationID == "" {
		err := sdkerrs.ErrArgs.Wrap("call api operationID is empty")
		log.ZError(ctx, "ApiRequest", err, "type", "ctx not set operationID")
		return err
	}
	defer func(start time.Time) {
		if err == nil {
			log.ZDebug(ctx, "CallApi", "api", api, "state", "success", "cost time", time.Since(start).Milliseconds())
		} else {
			log.ZError(ctx, "CallApi", err, "api", api, "state", "failed", "cost time", time.Since(start).Milliseconds())
		}
	}(time.Now())
	reqBody, err := json.Marshal(req)
	if err != nil {
		log.ZError(ctx, "ApiRequest", err, "type", "json.Marshal(req) failed")
		return sdkerrs.ErrSdkInternal.Wrap("json.Marshal(req) failed " + err.Error())
	}
	ctxInfo := ccontext.Info(ctx)
	reqUrl := ctxInfo.ApiAddr() + api
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(reqBody))
	if err != nil {
		log.ZError(ctx, "ApiRequest", err, "type", "http.NewRequestWithContext failed")
		return sdkerrs.ErrSdkInternal.Wrap("sdk http.NewRequestWithContext failed " + err.Error())
	}
	log.ZDebug(ctx, "ApiRequest", "url", reqUrl, "body", string(reqBody))
	request.ContentLength = int64(len(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("operationID", operationID)
	request.Header.Set("token", ctxInfo.Token())
	response, err := new(http.Client).Do(request)
	if err != nil {
		log.ZError(ctx, "ApiRequest", err, "type", "network error")
		return sdkerrs.ErrNetwork.Wrap("ApiPost http.Client.Do failed " + err.Error())
	}
	defer response.Body.Close()
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.ZError(ctx, "ApiResponse", err, "type", "read body", "status", response.Status)
		return sdkerrs.ErrSdkInternal.Wrap("io.ReadAll(ApiResponse) failed " + err.Error())
	}
	log.ZDebug(ctx, "ApiResponse", "url", reqUrl, "status", response.Status, "body", string(respBody))
	var baseApi apiResponse
	if err := json.Unmarshal(respBody, &baseApi); err != nil {
		log.ZError(ctx, "ApiResponse", err, "type", "api code parse")
		return sdkerrs.ErrSdkInternal.Wrap(fmt.Sprintf("api %s json.Unmarshal(%q, %T) failed %s", api, string(respBody), &baseApi, err.Error()))
	}
	if baseApi.ErrCode != 0 {
		err := sdkerrs.New(baseApi.ErrCode, baseApi.ErrMsg, baseApi.ErrDlt)
		log.ZError(ctx, "ApiResponse", err, "type", "api code error", "msg", baseApi.ErrMsg, "dlt", baseApi.ErrDlt)
		return err
	}
	if resp == nil || len(baseApi.Data) == 0 || string(baseApi.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(baseApi.Data, resp); err != nil {
		log.ZError(ctx, "ApiResponse", err, "type", "api data parse", "data", string(baseApi.Data), "bind", fmt.Sprintf("%T", resp))
		return sdkerrs.ErrSdkInternal.Wrap(fmt.Sprintf("json.Unmarshal(%q, %T) failed %s", string(baseApi.Data), resp, err.Error()))
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
	if req.GetPagination().ShowNumber <= 0 {
		req.GetPagination().ShowNumber = 50
	}
	var res []C
	for i := int32(0); ; i++ {
		req.GetPagination().PageNumber = i + 1
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
