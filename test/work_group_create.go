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
	"encoding/json"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/network"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"sync"
	"time"

	"github.com/OpenIMSDK/protocol/sdkws"
)

var (
	INVITEUSERTOGROUP = ""
)

func InviteListToGroup(userIDList []string, groupID string) {
	var inviteReq server_api_params.InviteUserToGroupReq
	inviteReq.OperationID = utils.OperationIDGenerator()
	inviteReq.GroupID = groupID
	inviteReq.InvitedUserIDList = userIDList
	for {
		resp, err := network.Post2Api(INVITEUSERTOGROUP, inviteReq, AdminToken)
		if err != nil {
			log.Warn(inviteReq.OperationID, " INVITE USER TO GROUP failed", inviteReq, "err: ", err)
			continue
		} else {
			log.Info(inviteReq.OperationID, " invite resp : ", string(resp))
			return
		}
	}
}

func InviteToGroup(userID string, groupID string) {
	var inviteReq server_api_params.InviteUserToGroupReq
	inviteReq.OperationID = utils.OperationIDGenerator()
	inviteReq.GroupID = groupID
	inviteReq.InvitedUserIDList = []string{userID}
	for {
		resp, err := network.Post2Api(INVITEUSERTOGROUP, inviteReq, AdminToken)
		if err != nil {
			log.Warn(inviteReq.OperationID, " INVITE USER TO GROUP failed", inviteReq, "err: ", err)
			continue
		} else {
			log.Info(inviteReq.OperationID, " invite resp : ", string(resp))
			return
		}
	}
}

func CreateWorkGroup(number int) string {
	t1 := time.Now()
	RegisterWorkGroupAccounts(number)
	log.Info("", "RegisterAccounts  cost time: ", time.Since(t1), "Online client number ", number)

	groupID := ""
	var req server_api_params.CreateGroupReq

	//var memberList []*server_api_params.GroupAddMemberInfo
	//for _, v := range allUserID {
	//	memberList = append(memberList, &server_api_params.GroupAddMemberInfo{UserID: v, RoleLevel: 1})
	//}
	//	req.MemberList = memberList
	req.OwnerUserID = "openIM123456"
	for {
		req.OperationID = utils.OperationIDGenerator()
		req.GroupType = constant.WorkingGroup
		req.OperationID = utils.OperationIDGenerator()
		resp, err := network.Post2Api(CREATEGROUP, req, AdminToken)
		if err != nil {
			log.Warn(req.OperationID, "CREATE GROUP failed", string(resp), "err: ", err)
			continue
		} else {
			type CreateGroupResp struct {
				server_api_params.CommResp
				GroupInfo sdkws.GroupInfo `json:"data"`
			}

			var result CreateGroupResp
			err := json.Unmarshal(resp, &result)
			if err != nil {
				log.Error(req.OperationID, "Unmarshal failed ", err.Error(), string(resp))
			}
			log.Info(req.OperationID, "Unmarshal  ", string(resp), result)
			groupID = result.GroupInfo.GroupID
			log.Info(req.OperationID, "create groupID:", groupID)
			break
		}
	}

	split := 100
	idx := 0
	remain := len(allUserID) % split
	for idx = 0; idx < len(allUserID)/split; idx++ {
		sub := allUserID[idx*split : (idx+1)*split]
		log.Warn(req.OperationID, "invite to groupID:", groupID)
		InviteListToGroup(sub, groupID)
	}
	if remain > 0 {
		sub := allUserID[idx*split:]
		log.Warn(req.OperationID, "invite to groupID:", groupID)
		InviteListToGroup(sub, groupID)
	}

	//var wg sync.WaitGroup
	//for _, v := range allUserID {
	//	wg.Add(1)
	//	go funcation(uID, gID string) {
	//		InviteToGroup(uID, gID)
	//		wg.Done()
	//	}(v, groupID)
	//}
	//wg.Wait()
	return groupID
}

func RegisterWorkGroupAccounts(number int) {
	var wg sync.WaitGroup
	wg.Add(number)
	for i := 0; i < number; i++ {
		go func(t int) {
			userID := GenUid(t, "workgroup")
			register(userID)
			log.Info("register ", userID)
			wg.Done()
		}(i)
	}
	wg.Wait()

	log.Info("", "RegisterAccounts finish ", number)
}

func RegisterWorkGroupPressAccounts(number int) {
	var wg sync.WaitGroup
	wg.Add(number)
	for i := 0; i < number; i++ {
		go func(t int) {
			userID := GenUid(t, "press_workgroup")
			register(userID)
			log.Info("register ", userID)
			wg.Done()
		}(i)
	}
	wg.Wait()

	userID1 := GenUid(1234567, "workgroup")
	register(userID1)
	userID2 := GenUid(7654321, "workgroup")
	register(userID2)
	log.Info("", "RegisterAccounts finish ", number)
}
