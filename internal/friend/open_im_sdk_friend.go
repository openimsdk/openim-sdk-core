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
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
)

//f
func (f *Friend) GetDesignatedFriendsInfo(callback open_im_sdk_callback.Base, friendUserIDList string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", friendUserIDList)
		var unmarshalList sdk.GetDesignatedFriendsInfoParams
		common.JsonUnmarshalCallback(friendUserIDList, &unmarshalList, callback, operationID)
		result := f.getDesignatedFriendsInfo(callback, unmarshalList, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(result))
	}()
}

func (f *Friend) AddFriend(callback open_im_sdk_callback.Base, userIDReqMsg string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", userIDReqMsg)
		var unmarshalAddFriendParams sdk.AddFriendParams
		common.JsonUnmarshalAndArgsValidate(userIDReqMsg, &unmarshalAddFriendParams, callback, operationID)
		f.addFriend(callback, unmarshalAddFriendParams, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.AddFriendCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk.AddFriendCallback))
	}()
}

func (f *Friend) GetRecvFriendApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ")
		result := f.getRecvFriendApplicationList(callback, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(result))
	}()
}

func (f *Friend) GetSendFriendApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ")
		result := f.getSendFriendApplicationList(callback, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(result))
	}()
}

func (f *Friend) AcceptFriendApplication(callback open_im_sdk_callback.Base, userIDHandleMsg string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", userIDHandleMsg)
		var unmarshalParams sdk.ProcessFriendApplicationParams
		common.JsonUnmarshalAndArgsValidate(userIDHandleMsg, &unmarshalParams, callback, operationID)
		f.processFriendApplication(callback, unmarshalParams, 1, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.ProcessFriendApplicationCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk.ProcessFriendApplicationCallback))
	}()
}

func (f *Friend) RefuseFriendApplication(callback open_im_sdk_callback.Base, userIDHandleMsg string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", userIDHandleMsg)
		var unmarshalParams sdk.ProcessFriendApplicationParams
		common.JsonUnmarshalAndArgsValidate(userIDHandleMsg, &unmarshalParams, callback, operationID)
		f.processFriendApplication(callback, unmarshalParams, -1, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.ProcessFriendApplicationCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk.ProcessFriendApplicationCallback))
	}()
}

func (f *Friend) CheckFriend(callback open_im_sdk_callback.Base, friendUserIDList string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", friendUserIDList)
		var unmarshalParams sdk.CheckFriendParams
		common.JsonUnmarshalAndArgsValidate(friendUserIDList, &unmarshalParams, callback, operationID)
		result := f.checkFriend(callback, unmarshalParams, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(result))
	}()
}

func (f *Friend) DeleteFriend(callback open_im_sdk_callback.Base, friendUserID string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", friendUserID)
		f.deleteFriend(sdk.DeleteFriendParams(friendUserID), callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.DeleteFriendCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk.DeleteFriendCallback))
	}()
}

//f
func (f *Friend) GetFriendList(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ")
		result := f.getFriendList(callback, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(result))
	}()
}
func (f *Friend) SearchFriends(callback open_im_sdk_callback.Base, searchParam, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", searchParam)
		var unmarshalSearchFriendsInfoParam sdk.SearchFriendsParam
		common.JsonUnmarshalAndArgsValidate(searchParam, &unmarshalSearchFriendsInfoParam, callback, operationID)
		unmarshalSearchFriendsInfoParam.KeywordList = utils.TrimStringList(unmarshalSearchFriendsInfoParam.KeywordList)
		friendsInfoList := f.searchFriends(callback, unmarshalSearchFriendsInfoParam, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(friendsInfoList))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(friendsInfoList), len(friendsInfoList))

	}()
}

func (f *Friend) SetFriendRemark(callback open_im_sdk_callback.Base, userIDRemark string, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", userIDRemark)
		var unmarshalParams sdk.SetFriendRemarkParams
		common.JsonUnmarshalAndArgsValidate(userIDRemark, &unmarshalParams, callback, operationID)
		f.setFriendRemark(unmarshalParams, callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.SetFriendRemarkCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk.SetFriendRemarkCallback))
	}()
}

func (f *Friend) AddBlack(callback open_im_sdk_callback.Base, blackUserID, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", blackUserID)
		f.addBlack(callback, sdk.AddBlackParams(blackUserID), operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.AddBlackCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk.AddBlackCallback))
	}()
}

func (f *Friend) GetBlackList(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ")
		localBlackList, err := f.db.GetBlackListDB()
		common.CheckDBErrCallback(callback, err, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(localBlackList))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonStringDefault(localBlackList))
	}()
}

func (f *Friend) RemoveBlack(callback open_im_sdk_callback.Base, blackUserID, operationID string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", blackUserID)
		f.removeBlack(callback, sdk.RemoveBlackParams(blackUserID), operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk.RemoveBlackCallback))
		log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk.RemoveBlackCallback))
	}()
}

func (f *Friend) SetFriendListener(listener open_im_sdk_callback.OnFriendshipListener) {
	if listener == nil {
		return
	}
	f.friendListener = listener
}
