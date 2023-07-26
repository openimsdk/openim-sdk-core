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
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"time"

	"github.com/OpenIMSDK/protocol/sdkws"
	userPB "github.com/OpenIMSDK/protocol/user"
)

func RegisterOne(uid, nickname, faceurl string) (bool, error) {
	InitContext(uid)
	// config.UserID = uid
	res, err := checkUserAccount(uid)
	if err != nil {
		return false, err
	}
	if res != true {
		log.Error(uid, "user:["+uid+"] account register fail, maybe already registered")
		return false, nil
	}
	return registerUserAccount(uid, nickname, faceurl), err
}

// 批量注册
// 返回值：成功注册和失败注册的 uidList
func RegisterBatch(users []Users) ([]string, []string) {
	var successList, failList []string
	for i := 0; i < len(users); i++ {
		uid := users[i].Uid
		nickname := users[i].Nickname
		faceUrl := users[i].FaceUrl
		bool, err := RegisterOne(uid, nickname, faceUrl)
		log.Info("", "注册：", bool, " uid:", uid)
		if err != nil {
			log.Error(utils.OperationIDGenerator(), err)
		}
		if bool == false {
			failList = append(failList, uid)
		} else {
			successList = append(successList, uid)
		}
	}
	return successList, failList
}

func checkUserAccount(uid string) (bool, error) {
	var getAccountCheckReq userPB.AccountCheckReq
	var getAccountCheckResp userPB.AccountCheckResp
	getAccountCheckReq.CheckUserIDs = []string{uid}
	for {
		err := util.ApiPost(ctx, ACCOUNT_CHECK, &getAccountCheckReq, &getAccountCheckResp)
		if err != nil {
			return false, err
		}
		if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == "registered" {
			log.Warn(getAccountCheckReq.CheckUserIDs[0], "Already registered ", uid, getAccountCheckResp.Results)
			userLock.Lock()
			AllUserID = append(AllUserID, uid)
			userLock.Unlock()
			return false, nil
		} else if len(getAccountCheckResp.Results) == 1 && getAccountCheckResp.Results[0].AccountStatus == "unregistered" {
			log.Info(getAccountCheckReq.CheckUserIDs[0], "not registered ", uid, getAccountCheckResp.Results)
			break
		} else {
			log.Error(getAccountCheckReq.CheckUserIDs[0], " failed, continue ", err, REGISTERADDR, getAccountCheckReq.CheckUserIDs)
			continue
		}
	}
	return true, nil
}

func registerUserAccount(uid, nickname, faceurl string) bool {
	var req userPB.UserRegisterReq
	req.Users = []*sdkws.UserInfo{{UserID: uid, Nickname: nickname, FaceURL: faceurl}}
	req.Secret = Secret
	for {
		err := util.ApiPost(ctx, USER_REGISTER, &req, nil)
		if err != nil {
			log.Error("post failed ,continue ", err.Error(), REGISTERADDR)
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			log.Info("register ok ", REGISTERADDR)
			userLock.Lock()
			AllUserID = append(AllUserID, uid)
			userLock.Unlock()
			return true
		}
	}
	return false
}
