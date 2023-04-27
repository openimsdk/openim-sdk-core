// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package testv2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"open_im_sdk/internal/util"
	"open_im_sdk/open_im_sdk"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

var (
	ctx = mcontext.NewCtx(utils.GetFuncName(2) + ":test")
)

func init() {
	rand.Seed(time.Now().UnixNano())
	listner := &OnConnListener{}
	config := getConf(APIADDR, WSADDR)
	util.BaseURL = APIADDR
	isInit := open_im_sdk.InitSDK(listner, "test", string(GetResValue(json.Marshal(config))))
	if !isInit {
		panic("init sdk failed")
	}
	ctx := mcontext.NewCtx("testInitLogin")
	token := GetResValue(GetUserToken(ctx, UserID))
	util.Token = token
	if err := open_im_sdk.UserForSDK.Login(ctx, UserID, token); err != nil {
		panic(err)
	}
	fmt.Println(util.Token)
}

func GetUserToken(ctx context.Context, userID string) (string, error) {
	jsonReqData, err := json.Marshal(map[string]any{
		"userID":   userID,
		"platform": 1,
		"secret":   "openIM123",
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, APIADDR+"/auth/user_token", bytes.NewReader(jsonReqData))
	if err != nil {
		return "", err
	}
	req.Header.Set("operationID", ctx.Value("operationID").(string))
	client := http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	type Result struct {
		ErrCode int    `json:"errCode"`
		ErrMsg  string `json:"errMsg"`
		ErrDlt  string `json:"errDlt"`
		Data    struct {
			Token             string `json:"token"`
			ExpireTimeSeconds int    `json:"expireTimeSeconds"`
		} `json:"data"`
	}
	var result Result
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("errCode:%d, errMsg:%s, errDlt:%s", result.ErrCode, result.ErrMsg, result.ErrDlt)
	}
	return result.Data.Token, nil
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetResValue[T any](value T, err error) T {
	CheckErr(err)
	return value
}
