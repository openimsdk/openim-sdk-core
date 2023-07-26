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
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"

	authPB "github.com/OpenIMSDK/protocol/auth"
)

func GetToken(uid string) (string, int64) {
	InitContext(uid)
	config.Token = ""
	req := authPB.UserTokenReq{
		PlatformID: PlatformID,
		UserID:     uid,
		Secret:     Secret,
	}
	resp := authPB.UserTokenResp{}
	err := util.ApiPost(ctx, RPC_USER_TOKEN, &req, &resp)
	if err != nil {
		log.Error(req.UserID, "ApiPost failed ", err.Error(), TOKENADDR, req)
		return "", 0
	}
	config.Token = resp.Token
	log.Info(req.UserID, "get token: ", resp.Token, " expireTimeSeconds: ", resp.ExpireTimeSeconds)
	return resp.Token, resp.ExpireTimeSeconds
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

func CreateCtx(uid string) context.Context {
	operationID := utils.OperationIDGenerator()
	ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID:   uid,
		Token:    AllLoginMgr[uid].Token,
		IMConfig: Config,
	})
	ctx = ccontext.WithOperationID(ctx, operationID)
	ctx = ccontext.WithSendMessageCallback(ctx, TestSendMsgCallBack{
		OperationID: operationID,
	})
	return ctx
}
