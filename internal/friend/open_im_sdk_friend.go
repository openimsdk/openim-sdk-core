// Copyright 2021 OpenIM Corporation
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

package friend

import (
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
)

func (f *Friend) GetDesignatedFriendsInfo(callback common.Base, friendUserIDList string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", friendUserIDList)
		var unmarshalList sdk.GetDesignatedFriendsInfoParams
		common.JsonUnmarshal(friendUserIDList, &unmarshalList, callback, operationID)
		result := f.getDesignatedFriendsInfo(callback, unmarshalList, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(result), result)
	}()
}

func (f *Friend) AddFriend(callback common.Base, userIDReqMsg string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", userIDReqMsg)
		var unmarshalAddFriendParams sdk.AddFriendParams
		common.JsonUnmarshalAndArgsValidate(userIDReqMsg, &unmarshalAddFriendParams, callback, operationID)
		f.addFriend(callback, unmarshalAddFriendParams, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.AddFriendCallback))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(sdk.AddFriendCallback))
	}()
}

func (f *Friend) GetRecvFriendApplicationList(callback common.Base, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
		result := f.getRecvFriendApplicationList(callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(result))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(result))
	}()
}

func (f *Friend) GetSendFriendApplicationList(callback common.Base, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
		result := f.getSendFriendApplicationList(callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(result))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(result))
	}()
}

func (f *Friend) AcceptFriendApplication(callback common.Base, userIDHandleMsg string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", userIDHandleMsg)
		var unmarshalParams sdk.ProcessFriendApplicationParams
		common.JsonUnmarshalAndArgsValidate(userIDHandleMsg, &unmarshalParams, callback, operationID)
		f.processFriendApplication(callback, unmarshalParams, 1, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.ProcessFriendApplicationCallback))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(sdk.ProcessFriendApplicationCallback))
	}()
}

func (f *Friend) RefuseFriendApplication(callback common.Base, userIDHandleMsg string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", userIDHandleMsg)
		var unmarshalParams sdk.ProcessFriendApplicationParams
		common.JsonUnmarshalAndArgsValidate(userIDHandleMsg, &unmarshalParams, callback, operationID)
		f.processFriendApplication(callback, unmarshalParams, -1, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.ProcessFriendApplicationCallback))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(sdk.ProcessFriendApplicationCallback))
	}()
}

func (f *Friend) CheckFriend(callback common.Base, friendUserIDList string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", friendUserIDList)
		var unmarshalParams sdk.CheckFriendParams
		common.JsonUnmarshalAndArgsValidate(friendUserIDList, &unmarshalParams, callback, operationID)
		result := f.checkFriend(callback, unmarshalParams, operationID)
		callback.OnSuccess(utils.StructToJsonString(result))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(result))
	}()
}

func (f *Friend) DeleteFriend(callback common.Base, friendUserID string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", friendUserID)
		f.deleteFriend(sdk.DeleteFriendParams(friendUserID), callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.DeleteFriendCallback))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(sdk.DeleteFriendCallback))
	}()
}

func (f *Friend) GetFriendList(callback common.Base, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
		var filterLocalFriendList sdk.GetFriendListCallback
		localFriendList, err := f.db.GetAllFriendList()
		common.CheckErr(callback, err, operationID)
		localBlackUidList, err := f.db.GetBlackListUserID()
		common.CheckErr(callback, err, operationID)
		for _, v := range localFriendList {
			if !utils.IsContain(v.FriendUserID, localBlackUidList) {
				filterLocalFriendList = append(filterLocalFriendList, v)
			}
		}
		callback.OnSuccess(utils.StructToJsonString(filterLocalFriendList))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(filterLocalFriendList))
	}()
}

func (f *Friend) SetFriendRemark(callback common.Base, userIDRemark string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", userIDRemark)
		var unmarshalParams sdk.SetFriendRemarkParams
		common.JsonUnmarshalAndArgsValidate(userIDRemark, &unmarshalParams, callback, operationID)
		f.setFriendRemark(unmarshalParams, callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.SetFriendRemarkCallback))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(sdk.SetFriendRemarkCallback))
	}()
}

func (f *Friend) AddBlack(callback common.Base, blackUserID, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", blackUserID)
		f.addBlack(callback, sdk.AddBlackParams(blackUserID), operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.AddBlackCallback))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(sdk.AddBlackCallback))
	}()
}

func (f *Friend) GetBlackList(callback common.Base, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
		localBlackList, err := f.db.GetBlackList()
		common.CheckErr(callback, err, operationID)
		callback.OnSuccess(utils.StructToJsonString(localBlackList))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(localBlackList))
	}()
}

func (f *Friend) RemoveBlack(callback common.Base, blackUserID, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", blackUserID)
		f.removeBlack(callback, sdk.RemoveBlackParams(blackUserID), operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.RemoveBlackCallback))
		log.NewInfo(operationID, utils.GetSelfFuncName(), " callback: ", utils.StructToJsonString(sdk.RemoveBlackCallback))
	}()
}

func (f *Friend) SetFriendListener(listener OnFriendshipListener) {
	if listener == nil {
		return
	}
	f.friendListener = listener
}
