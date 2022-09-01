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
)

func CreateWorkGroup(number int) string {
	t1 := time.Now()
	RegisterWorkGroupAccounts(number)
	log.Info("", "RegisterAccounts  cost time: ", time.Since(t1), "Online client number ", number)

	var req server_api_params.CreateGroupReq
	req.OperationID = utils.OperationIDGenerator()
	req.GroupType = constant.WorkingGroup
	req.OperationID = utils.OperationIDGenerator()
	var memberList []*server_api_params.GroupAddMemberInfo
	for _, v := range allUserID {
		memberList = append(memberList, &server_api_params.GroupAddMemberInfo{UserID: v, RoleLevel: 1})
	}
	req.MemberList = memberList
	req.OwnerUserID = "openIM123456"
	for {
		resp, err := network.Post2Api(CREATEGROUP, req, AdminToken)
		if err != nil {
			log.Warn(req.OperationID, "CREATE GROUP failed", string(resp), "err: ", err)
			continue
		} else {
			var result server_api_params.CreateGroupResp
			json.Unmarshal(resp, result)
			return result.GroupInfo.GroupID
		}
	}
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
