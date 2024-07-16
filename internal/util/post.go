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
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/page"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

// apiClient is a global HTTP client with a timeout of one minute.
var apiClient = &http.Client{
	Timeout: time.Second * 30,
}

// ApiResponse represents the standard structure of an API response.
type ApiResponse struct {
	ErrCode int             `json:"errCode"`
	ErrMsg  string          `json:"errMsg"`
	ErrDlt  string          `json:"errDlt"`
	Data    json.RawMessage `json:"data"`
}

// ApiPost performs an HTTP POST request to a specified API endpoint.
// It serializes the request object, sends it to the API, and unmarshals the response into the resp object.
// It handles logging, error wrapping, and operation ID validation.
// Context (ctx) is used for passing metadata and control information.
// api: the API endpoint to which the request is sent.
// req: the request object to be sent to the API.
// resp: a pointer to the response object where the API response will be unmarshalled.
// Returns an error if the request fails at any stage.
func ApiPost(ctx context.Context, api string, req, resp any) (err error) {
	// Extract operationID from context and validate.
	operationID, _ := ctx.Value("operationID").(string)
	if operationID == "" {
		err := sdkerrs.ErrArgs.WrapMsg("call api operationID is empty")
		log.ZError(ctx, "ApiRequest", err, "type", "ctx not set operationID")
		return err
	}

	// Deferred function to log the result of the API call.
	defer func(start time.Time) {
		elapsed := time.Since(start).Milliseconds()
		if err == nil {
			log.ZDebug(ctx, "CallApi", "api", api, "state", "success", "cost time", fmt.Sprintf("%dms", elapsed))
		} else {
			log.ZError(ctx, "CallApi", err, "api", api, "state", "failed", "cost time", fmt.Sprintf("%dms", elapsed))
		}
	}(time.Now())

	// Serialize the request object to JSON.
	reqBody, err := json.Marshal(req)
	if err != nil {
		log.ZError(ctx, "ApiRequest", err, "type", "json.Marshal(req) failed")
		return sdkerrs.ErrSdkInternal.WrapMsg("json.Marshal(req) failed " + err.Error())
	}

	// Construct the full API URL and create a new HTTP request with context.
	ctxInfo := ccontext.Info(ctx)
	reqUrl := ctxInfo.ApiAddr() + api
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewReader(reqBody))
	if err != nil {
		log.ZError(ctx, "ApiRequest", err, "type", "http.NewRequestWithContext failed")
		return sdkerrs.ErrSdkInternal.WrapMsg("sdk http.NewRequestWithContext failed " + err.Error())
	}

	// Set headers for the request.
	log.ZDebug(ctx, "ApiRequest", "url", reqUrl, "body", string(reqBody))
	request.ContentLength = int64(len(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("operationID", operationID)
	request.Header.Set("token", ctxInfo.Token())

	// Send the request and receive the response.
	response, err := apiClient.Do(request)
	if err != nil {
		log.ZError(ctx, "ApiRequest", err, "type", "network error")
		return sdkerrs.ErrNetwork.WrapMsg("ApiPost http.Client.Do failed " + err.Error())
	}

	// Ensure the response body is closed after processing.
	defer response.Body.Close()

	// Read the response body.
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.ZError(ctx, "ApiResponse", err, "type", "read body", "status", response.Status)
		return sdkerrs.ErrSdkInternal.WrapMsg("io.ReadAll(ApiResponse) failed " + err.Error())
	}

	// Log the response for debugging purposes.
	log.ZDebug(ctx, "ApiResponse", "url", reqUrl, "status", response.Status, "body", string(respBody))

	// Unmarshal the response body into the ApiResponse structure.
	var baseApi ApiResponse
	if err := json.Unmarshal(respBody, &baseApi); err != nil {
		log.ZError(ctx, "ApiResponse", err, "type", "api code parse")
		return sdkerrs.ErrSdkInternal.WrapMsg(fmt.Sprintf("api %s json.Unmarshal(%q, %T) failed %s", api, string(respBody), &baseApi, err.Error()))
	}

	// Check if the API returned an error code and handle it.
	if baseApi.ErrCode != 0 {
		err := sdkerrs.New(baseApi.ErrCode, baseApi.ErrMsg, baseApi.ErrDlt)
		ccontext.GetApiErrCodeCallback(ctx).OnError(ctx, err)
		log.ZError(ctx, "ApiResponse", err, "type", "api code error", "msg", baseApi.ErrMsg, "dlt", baseApi.ErrDlt)
		return err
	}

	// If no data is received, or it's null, return with no error.
	if resp == nil || len(baseApi.Data) == 0 || string(baseApi.Data) == "null" {
		return nil
	}

	// Unmarshal the actual data part of the response into the provided response object.
	if err := json.Unmarshal(baseApi.Data, resp); err != nil {
		log.ZError(ctx, "ApiResponse", err, "type", "api data parse", "data", string(baseApi.Data), "bind", fmt.Sprintf("%T", resp))
		return sdkerrs.ErrSdkInternal.WrapMsg(fmt.Sprintf("json.Unmarshal(%q, %T) failed %s", string(baseApi.Data), resp, err.Error()))
	}

	return nil
}

// CallApi wraps ApiPost to make an API call and unmarshal the response into a new instance of type T.
func CallApi[T any](ctx context.Context, api string, req any) (*T, error) {
	var resp T
	if err := ApiPost(ctx, api, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPageAll handles pagination for API requests. It iterates over pages of data until all data is retrieved.
// A is the request type with pagination support, B is the response type, and C is the type of data to be returned.
// The function fn processes each page of response data to extract a slice of C.
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

func GetPageAllWithMaxNum[A interface {
	GetPagination() *sdkws.RequestPagination
}, B, C any](ctx context.Context, api string, req A, fn func(resp *B) []C, maxItems int) ([]C, error) {
	if req.GetPagination().ShowNumber <= 0 {
		req.GetPagination().ShowNumber = 50
	}
	var res []C
	totalFetched := 0
	for i := int32(0); ; i++ {
		req.GetPagination().PageNumber = i + 1
		memberResp, err := CallApi[B](ctx, api, req)
		if err != nil {
			return nil, err
		}
		list := fn(memberResp)
		res = append(res, list...)
		totalFetched += len(list)
		if len(list) < int(req.GetPagination().ShowNumber) || (maxItems > 0 && totalFetched >= maxItems) {
			break
		}
	}
	if maxItems > 0 && len(res) > maxItems {
		res = res[:maxItems]
	}
	return res, nil
}

func FetchAndInsertPagedData[RESP, L any](ctx context.Context, api string, req page.PageReq, fn func(resp *RESP) []L, batchInsertFn func(ctx context.Context, items []L) error,
	insertFn func(ctx context.Context, item L) error, maxItems int64) error {
	if req.GetPagination().ShowNumber <= 0 {
		req.GetPagination().ShowNumber = 50
	}
	var errSingle error
	var errList []error
	totalFetched := 0
	for i := int32(0); ; i++ {
		req.GetPagination().PageNumber = i + 1
		memberResp, err := CallApi[RESP](ctx, api, req)
		if err != nil {
			return err
		}
		list := fn(memberResp)
		if err := batchInsertFn(ctx, list); err != nil {
			for _, item := range list {
				errSingle = insertFn(ctx, item)
				if errSingle != nil {
					errList = append(errList, errs.New(errSingle.Error(), "item", item))
				}
			}
		}
		totalFetched += len(list)
		if len(list) < int(req.GetPagination().ShowNumber) || (maxItems > 0 && totalFetched >= int(maxItems)) {
			break
		}
	}
	if len(errList) > 0 {
		return errs.WrapMsg(errList[0], "batch insert failed due to data exception")
	}
	return nil
}
