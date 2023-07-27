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

package test

import (
	"context"
	"fmt"
	"net"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/network"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"os"
	"strconv"
	"sync"
	"time"

	authPB "github.com/OpenIMSDK/protocol/auth"
	"github.com/OpenIMSDK/protocol/sdkws"
	userPB "github.com/OpenIMSDK/protocol/user"
	imLog "github.com/OpenIMSDK/tools/log"
)

func GenUid(uid int, prefix string) string {
	if getMyIP() == "" {
		log.Error("", "getMyIP() failed, exit ")
		os.Exit(1)
	}
	UidPrefix := getMyIP() + "_" + prefix + "_"
	return UidPrefix + strconv.FormatInt(int64(uid), 10)
}

func RegisterOnlineAccounts(number int) {
	var wg sync.WaitGroup
	wg.Add(number)
	for i := 0; i < number; i++ {
		go func(t int) {
			userID := GenUid(t, "online")
			register(userID)
			log.Info("register ", userID)
			wg.Done()
		}(i)

	}
	wg.Wait()
	log.Info("", "RegisterAccounts finish ", number)
}

type GetTokenReq struct {
	Secret   string `json:"secret"`
	Platform int    `json:"platform"`
	Uid      string `json:"uid"`
}

type RegisterReq struct {
	Secret   string `json:"secret"`
	Platform int    `json:"platform"`
	Uid      string `json:"uid"`
	Name     string `json:"name"`
}

type ResToken struct {
	Data struct {
		ExpiredTime int64  `json:"expiredTime"`
		Token       string `json:"token"`
		Uid         string `json:"uid"`
	}
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

var AdminToken = ""

func init() {
	AdminToken = getToken("openIM123456")
	if err := imLog.InitFromConfig("open-im-sdk-core", "", int(LogLevel), IsLogStandardOutput, false, LogFilePath, 0, 24); err != nil {
		fmt.Println("123456", "log init failed ", err.Error())
	}
}

var ctx context.Context
var cctx context.Context

func register(uid string) error {
	ctx = ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: uid,
		Token:  AdminToken,
		IMConfig: sdk_struct.IMConfig{
			PlatformID: PlatformID,
			ApiAddr:    APIADDR,
			WsAddr:     WSADDR,
			LogLevel:   LogLevel,
		},
	})
	ctx = ccontext.WithOperationID(ctx, "123456")

	//ACCOUNTCHECK
	var getAccountCheckReq userPB.AccountCheckReq
	var getAccountCheckResp userPB.AccountCheckResp
	getAccountCheckReq.CheckUserIDs = []string{uid}

	for {
		err := util.ApiPost(ctx, "/user/account_check", &getAccountCheckReq, &getAccountCheckResp)
		if err != nil {
			return err
		}
		if len(getAccountCheckResp.Results) == 1 &&
			getAccountCheckResp.Results[0].AccountStatus == "registered" {
			log.Warn(getAccountCheckReq.CheckUserIDs[0], "Already registered ", uid, getAccountCheckResp)
			userLock.Lock()
			allUserID = append(allUserID, uid)
			userLock.Unlock()
			return nil
		} else if len(getAccountCheckResp.Results) == 1 &&
			getAccountCheckResp.Results[0].AccountStatus == "unregistered" {
			log.Info(getAccountCheckReq.CheckUserIDs[0], "not registered ", uid, getAccountCheckResp)
			break
		} else {
			log.Error(getAccountCheckReq.CheckUserIDs[0], " failed, continue ", err, REGISTERADDR, getAccountCheckReq)
			continue
		}
	}

	var rreq userPB.UserRegisterReq
	rreq.Users = []*sdkws.UserInfo{{UserID: uid}}

	for {
		err := util.ApiPost(ctx, "/auth/user_register", &rreq, nil)
		if err != nil {
			log.Error("post failed ,continue ", err.Error(), REGISTERADDR, getAccountCheckReq)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.Info("register ok ", REGISTERADDR, getAccountCheckReq)
			userLock.Lock()
			allUserID = append(allUserID, uid)
			userLock.Unlock()
			return nil
		}
	}
}

func getToken(uid string) string {
	ctx = ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: uid,
		Token:  "",
		IMConfig: sdk_struct.IMConfig{
			PlatformID: PlatformID,
			ApiAddr:    APIADDR,
			WsAddr:     WSADDR,
			LogLevel:   LogLevel,
		},
	})
	ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
	url := TOKENADDR
	req := authPB.UserTokenReq{
		Secret:     SECRET,
		PlatformID: PlatformID,
		UserID:     uid,
	}
	resp := authPB.UserTokenResp{}
	err := util.ApiPost(ctx, "/auth/user_token", &req, &resp)
	if err != nil {
		log.Error(req.UserID, "Post2Api failed ", err.Error(), url, req)
		return ""
	}

	log.Info(req.UserID, "get token: ", resp.Token)
	return resp.Token
}

func RunGetToken(strMyUid string) string {
	var token string
	for true {
		token = getToken(strMyUid)
		if token == "" {
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		} else {
			break
		}
	}
	return token
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

func RegisterReliabilityUser(id int, timeStamp string) {
	userID := GenUid(id, "reliability_"+timeStamp)
	register(userID)
	token := RunGetToken(userID)
	coreMgrLock.Lock()
	defer coreMgrLock.Unlock()
	allLoginMgr[id] = &CoreNode{token: token, userID: userID}
}

func WorkGroupRegisterReliabilityUser(id int) {
	userID := GenUid(id, "workgroup")
	//	register(userID)
	token := RunGetToken(userID)
	coreMgrLock.Lock()
	defer coreMgrLock.Unlock()
	log.Info("", "WorkGroupRegisterReliabilityUser userID: ", userID, "token: ", token)
	allLoginMgr[id] = &CoreNode{token: token, userID: userID}
}

func RegisterPressUser(id int) {
	userID := GenUid(id, "press")
	register(userID)
	token := RunGetToken(userID)
	coreMgrLock.Lock()
	defer coreMgrLock.Unlock()
	allLoginMgr[id] = &CoreNode{token: token, userID: userID}
}

func GetGroupMemberNum(groupID string) uint32 {
	var req server_api_params.GetGroupInfoReq
	req.OperationID = utils.OperationIDGenerator()
	req.GroupIDList = []string{groupID}

	var groupInfoList []*sdkws.GroupInfo

	r, err := network.Post2Api(GETGROUPSINFOROUTER, req, AdminToken)
	if err != nil {
		log.Error("", "post failed ", GETGROUPSINFOROUTER, req)
		return 0
	}
	err = common.CheckErrAndResp(nil, r, &groupInfoList, nil)
	if err != nil {
		log.Error("", "CheckErrAndResp failed ", err.Error(), string(r))
		return 0
	}
	log.Warn("", "group info", groupInfoList)
	return groupInfoList[0].MemberCount
}
