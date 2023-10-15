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

package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"io/ioutil"
	"net/http"
	"time"
)

func get(url string) (response []byte, err error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func retry(ctx context.Context, url string, data interface{}, token string, attempts int, sleep time.Duration, fn func(context.Context, string, interface{}, string) ([]byte, error)) ([]byte, error) {
	if ctx.Err() == context.Canceled {
		return nil, errors.New("retry aborted due to context cancellation")
	}
	b, err := fn(ctx, url, data, token)
	if err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(ctx, url, data, token, attempts, 2*sleep, fn)
		}
		return nil, err
	}
	return b, nil
}

// application/json; charset=utf-8
func Post2Api(ctx context.Context, url string, data interface{}, token string) (content []byte, err error) {
	if ctx.Err() == context.Canceled {
		return nil, errors.New("Post2Api aborted due to context cancellation")
	}
	c, err := postLogic(ctx, url, data, token)
	return c, utils.Wrap(err, " post")
	return retry(ctx, url, data, token, 1, 10*time.Second, postLogic)
}

func Post2ApiForRead(ctx context.Context, url string, data interface{}, token string) (content []byte, err error) {
	if ctx.Err() == context.Canceled {
		return nil, errors.New("Post2ApiForRead aborted due to context cancellation")
	}
	return retry(ctx, url, data, token, 3, 10*time.Second, postLogic)
}

func postLogic(ctx context.Context, url string, data interface{}, token string) (content []byte, err error) {
	if ctx.Err() == context.Canceled {
		return nil, errors.New("postLogic aborted due to context cancellation")
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, utils.Wrap(err, "marshal failed, url")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, utils.Wrap(err, "newRequest failed, url")
	}
	req.Close = true
	req.Header.Add("content-type", "application/json")
	req.Header.Add("token", token)
	req.Header.Add("OperationID", utils.OperationIDGenerator())
	tp := &http.Transport{
		DialContext: (&net.Dialer{
			KeepAlive: 10 * time.Minute,
		}).DialContext,
		ResponseHeaderTimeout: 60 * time.Second,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
	}
	client := &http.Client{Timeout: 60 * time.Second, Transport: tp}
	resp, err := client.Do(req)
	if err != nil {
		return nil, utils.Wrap(err, "client.Do failed, url")
	}
	if resp.StatusCode != 200 {
		return nil, utils.Wrap(errors.New(resp.Status), "status code failed "+url)
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.Wrap(err, "ioutil.ReadAll failed, url")
	}
	//	fmt.Println(url, "Marshal data: ", string(jsonStr), string(result))
	return result, nil
}
