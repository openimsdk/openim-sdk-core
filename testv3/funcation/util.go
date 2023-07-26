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

package funcation

import (
	"context"
	"net"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"os"
	"strconv"
	"time"

	authPB "github.com/OpenIMSDK/protocol/auth"
)

func GenUid(uid int, prefix string) string {
	if getMyIP() == "" {
		log.Error("", "getMyIP() failed, exit ")
		os.Exit(1)
	}
	UidPrefix := getMyIP() + "_" + prefix + "_"
	return UidPrefix + strconv.FormatInt(int64(uid), 10)
}

func getMyIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error("", "InterfaceAddrs failed ", err.Error())
		os.Exit(1)
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func getToken(uid string) (string, int64) {
	InitContext(uid)
	config.Token = ""
	req := authPB.UserTokenReq{PlatformID: PlatformID, UserID: uid, Secret: Secret}
	resp := authPB.UserTokenResp{}
	err := util.ApiPost(ctx, "/auth/user_token", &req, &resp)
	if err != nil {
		log.Error(req.UserID, "ApiPost failed ", err.Error(), TOKENADDR, req)
		return "", 0
	}
	config.Token = resp.Token
	log.Info(req.UserID, "get token: ", resp.Token, " expireTimeSeconds: ", resp.ExpireTimeSeconds)
	return resp.Token, resp.ExpireTimeSeconds
}

func RunGetToken(strMyUid string) (string, int64) {
	var token string
	var exprie int64
	for true {
		token, exprie = getToken(strMyUid)
		if token == "" {
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		} else {
			break
		}
	}
	return token, exprie
}

func InitContext(uid string) context.Context {
	config = ccontext.GlobalConfig{
		UserID: uid,
		Token:  AdminToken,
		IMConfig: sdk_struct.IMConfig{
			PlatformID: PlatformID,
			ApiAddr:    APIADDR,
			WsAddr:     WSADDR,
			LogLevel:   LogLevel,
		},
	}
	ctx = ccontext.WithInfo(context.Background(), &config)
	ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
	return ctx
}
