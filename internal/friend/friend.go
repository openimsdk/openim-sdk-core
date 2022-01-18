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
	"github.com/mitchellh/mapstructure"
	comm "open_im_sdk/internal/common"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"strings"
)

type Friend struct {
	friendListener OnFriendshipListener
	loginUserID    string
	db             *db.DataBase
	p              *ws.PostApi
}

type OnFriendshipListener interface {
	OnFriendApplicationAdded(friendApplication string)
	OnFriendApplicationDeleted(friendApplication string)
	OnFriendApplicationAccepted(groupApplication string)
	OnFriendApplicationRejected(friendApplication string)
	OnFriendAdded(friendInfo string)
	OnFriendDeleted(friendInfo string)
	OnFriendInfoChanged(friendInfo string)
	OnBlackAdded(blackInfo string)
	OnBlackDeleted(blackInfo string)
}

func NewFriend(loginUserID string, db *db.DataBase, p *ws.PostApi) *Friend {
	return &Friend{loginUserID: loginUserID, db: db, p: p}
}

func (f *Friend) SetListener(listener OnFriendshipListener) {
	f.friendListener = listener
}

func (f *Friend) getDesignatedFriendsInfo(callback common.Base, friendUserIDList sdk.GetDesignatedFriendsInfoParams, operationID string) sdk.GetDesignatedFriendsInfoCallback {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", friendUserIDList)
	blackList, err := f.db.GetBlackInfoList(friendUserIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	var pureFriendUserIDList []string
	for _, v := range friendUserIDList {
		flag := 0
		for _, k := range blackList {
			if v == k.BlockUserID {
				flag = 1
				break
			}
		}
		if flag == 0 {
			pureFriendUserIDList = append(pureFriendUserIDList, v)
		}
	}
	localFriendList, err := f.db.GetFriendInfoList(pureFriendUserIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", pureFriendUserIDList, localFriendList)
	return localFriendList
}

func (f *Friend) addFriend(callback common.Base, userIDReqMsg sdk.AddFriendParams, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", userIDReqMsg)
	apiReq := api.AddFriendReq{}
	apiReq.ToUserID = userIDReqMsg.ToUserID
	apiReq.FromUserID = f.loginUserID
	apiReq.ReqMsg = userIDReqMsg.ReqMsg
	apiReq.OperationID = operationID
	f.p.PostFatalCallback(callback, constant.AddFriendRouter, apiReq, nil, operationID)
}

func (f *Friend) getRecvFriendApplicationList(callback common.Base, operationID string) sdk.GetRecvFriendApplicationListCallback {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	friendApplicationList, err := f.db.GetRecvFriendApplication()
	common.CheckDBErrCallback(callback, err, operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", friendApplicationList)
	return friendApplicationList
}

func (f *Friend) getSendFriendApplicationList(callback common.Base, operationID string) sdk.GetSendFriendApplicationListCallback {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	friendApplicationList, err := f.db.GetSendFriendApplication()
	common.CheckDBErrCallback(callback, err, operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", friendApplicationList)
	return friendApplicationList
}

func (f *Friend) processFriendApplication(callback common.Base, userIDHandleMsg sdk.ProcessFriendApplicationParams, handleResult int32, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", userIDHandleMsg, handleResult)
	apiReq := api.AddFriendResponseReq{}
	apiReq.FromUserID = f.loginUserID
	apiReq.ToUserID = userIDHandleMsg.ToUserID
	apiReq.Flag = handleResult
	apiReq.OperationID = operationID
	apiReq.HandleMsg = userIDHandleMsg.HandleMsg
	f.p.PostFatalCallback(callback, constant.AddFriendResponse, apiReq, nil, operationID)
	f.SyncFriendApplication(operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ")
}

func (f *Friend) checkFriend(callback common.Base, friendUserIDList sdk.CheckFriendParams, operationID string) sdk.CheckFriendCallback {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", friendUserIDList)
	friendList, err := f.db.GetFriendInfoList(friendUserIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	blackList, err := f.db.GetBlackInfoList(friendUserIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	var checkFriendCallback sdk.CheckFriendCallback
	for _, v := range friendUserIDList {
		var r api.UserIDResult
		isBlack := false
		isFriend := false
		for _, b := range blackList {
			if v == b.BlockUserID {
				isBlack = true
				break
			}
		}
		for _, f := range friendList {
			if v == f.FriendUserID {
				isFriend = true
				break
			}
		}
		r.UserID = v
		if isFriend && !isBlack {
			r.Result = 1
		} else {
			r.Result = 0
		}
		checkFriendCallback = append(checkFriendCallback, r)
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", checkFriendCallback)
	return checkFriendCallback
}

func (f *Friend) deleteFriend(friendUserID sdk.DeleteFriendParams, callback common.Base, operationID string) *api.CommDataResp {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", friendUserID)
	apiReq := api.DeleteFriendReq{}
	apiReq.ToUserID = string(friendUserID)
	apiReq.FromUserID = f.loginUserID
	apiReq.OperationID = operationID
	result := f.p.PostFatalCallback(callback, constant.DeleteFriendRouter, apiReq, operationID)
	f.SyncFriendList(operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", result)
	return result
}

func (f *Friend) setFriendRemark(userIDRemark sdk.SetFriendRemarkParams, callback common.Base, operationID string) *api.CommDataResp {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", userIDRemark)
	apiReq := api.SetFriendRemarkReq{}
	apiReq.OperationID = operationID
	apiReq.ToUserID = userIDRemark.ToUserID
	apiReq.FromUserID = f.loginUserID
	result := f.p.PostFatalCallback(callback, constant.SetFriendRemark, apiReq, operationID)
	f.SyncFriendList(operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", result)
	return result
	//
	//c := ConversationStruct{
	//	ConversationID: GetConversationIDBySessionType(uid2comm.Uid, SingleChatType),
	//}
	//faceUrl, name, err := u.getUserNameAndFaceUrlByUid(uid2comm.Uid)
	//if err != nil {
	//	sdkLog("getUserNameAndFaceUrlByUid err:", err)
	//	return
	//}
	//c.FaceURL = faceUrl
	//c.ShowName = name
	//u.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, UpdateFaceUrlAndNickName, c}})
	//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{c.ConversationID}}})
}

func (f *Friend) getServerFriendList(operationID string) ([]*api.FriendInfo, error) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	apiReq := api.GetFriendListReq{OperationID: operationID, FromUserID: f.loginUserID}
	resp, err := f.p.PostReturn(constant.GetFriendListRouter, apiReq)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	realData := api.GetFriendListResp{}
	err = mapstructure.Decode(resp.Data, &realData.FriendInfoList)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", realData.FriendInfoList)
	return realData.FriendInfoList, nil
}

func (f *Friend) getServerBlackList(operationID string) ([]*api.PublicUserInfo, error) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	apiReq := api.GetBlackListReq{OperationID: operationID, FromUserID: f.loginUserID}
	commData, err := f.p.PostReturn(constant.GetBlackListRouter, apiReq)
	if err != nil {
		return nil, utils.Wrap(err, operationID)
	}
	realData := api.GetBlackListResp{}
	err = mapstructure.Decode(commData.Data, &realData.BlackUserInfoList)
	if err != nil {
		return nil, utils.Wrap(err, operationID)
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", realData.BlackUserInfoList)
	return realData.BlackUserInfoList, nil
}

//recv
func (f *Friend) getFriendApplicationFromServer(operationID string) ([]*api.FriendRequest, error) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	apiReq := api.GetFriendApplyListReq{OperationID: operationID, FromUserID: f.loginUserID}
	resp, err := f.p.PostReturn(constant.GetFriendApplicationListRouter, apiReq)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	realData := api.GetFriendApplyListResp{}
	if err = mapstructure.Decode(resp.Data, &realData.FriendRequestList); err != nil {
		return nil, utils.Wrap(err, operationID)
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", realData.FriendRequestList)
	return realData.FriendRequestList, nil
}

//send
func (f *Friend) getSelfFriendApplicationFromServer(operationID string) ([]*api.FriendRequest, error) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	apiReq := api.GetSelfFriendApplyListReq{OperationID: operationID, FromUserID: f.loginUserID}
	resp, err := f.p.PostReturn(constant.GetSelfFriendApplicationListRouter, apiReq)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	realData := api.GetSelfFriendApplyListResp{}
	if err = mapstructure.Decode(resp.Data, &realData.FriendRequestList); err != nil {
		return nil, utils.Wrap(err, operationID)
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", realData.FriendRequestList)
	return realData.FriendRequestList, nil
}

func (f *Friend) addBlack(callback common.Base, blackUserID sdk.AddBlackParams, operationID string) *api.CommDataResp {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", blackUserID)
	apiReq := api.AddBlacklistReq{}
	apiReq.ToUserID = string(blackUserID)
	apiReq.FromUserID = f.loginUserID
	apiReq.OperationID = operationID
	result := f.p.PostFatalCallback(callback, constant.AddBlackRouter, apiReq, operationID)
	f.SyncBlackList(operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", result)
	return result
}

func (f *Friend) removeBlack(callback common.Base, blackUserID sdk.RemoveBlackParams, operationID string) *api.CommDataResp {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", blackUserID)
	apiReq := api.RemoveBlackListReq{}
	apiReq.ToUserID = string(blackUserID)
	apiReq.FromUserID = f.loginUserID
	apiReq.OperationID = operationID
	result := f.p.PostFatalCallback(callback, constant.RemoveBlackRouter, apiReq, operationID)
	f.SyncBlackList(operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", result)
	return result
}

func (f *Friend) SyncSelfFriendApplication(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := f.getSelfFriendApplicationFromServer(operationID)
	if err != nil {
		log.NewError(operationID, "getSelfFriendApplicationFromServer failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalFriendRequest(svrList)
	onLocal, err := f.db.GetSendFriendApplication()
	if err != nil {
		log.NewError(operationID, "GetSendFriendApplication failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "list", svrList, onServer, onLocal)

	aInBNot, bInANot, sameA, sameB := common.CheckFriendRequestDiff(onServer, onLocal)
	log.Debug(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := f.db.InsertFriendRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertFriendRequest failed ", err.Error())
			continue
		}
		callbackData := sdk.FriendApplicationAddedCallback(*onServer[index])
		f.friendListener.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
	}
	for _, index := range sameA {
		err := f.db.UpdateFriendRequest(onServer[index])
		if err != nil {
			if !strings.Contains(err.Error(), "RowsAffected == 0") {
				log.NewError(operationID, "UpdateFriendRequest failed ", err.Error(), *onServer[index])
				continue
			}
			if onServer[index].HandleResult == -1 {
				callbackData := sdk.FriendApplicationRejectCallback(*onServer[index])
				f.friendListener.OnFriendApplicationRejected(utils.StructToJsonString(callbackData))

			} else if onServer[index].HandleResult == 1 {
				callbackData := sdk.FriendApplicationAcceptCallback(*onServer[index])
				f.friendListener.OnFriendApplicationAccepted(utils.StructToJsonString(callbackData))
			}
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriendRequestBothUserID(onLocal[index].FromUserID, onLocal[index].ToUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriendRequestBothUserID failed ", err.Error(), onLocal[index].FromUserID, onLocal[index].ToUserID)
			continue
		}
		callbackData := sdk.FriendApplicationAcceptCallback(*onLocal[index])
		f.friendListener.OnFriendApplicationDeleted(utils.StructToJsonString(callbackData))
	}

}

func (f *Friend) SyncFriendApplication(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := f.getFriendApplicationFromServer(operationID)
	if err != nil {
		log.NewError(operationID, "getServerFriendList failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalFriendRequest(svrList)
	onLocal, err := f.db.GetRecvFriendApplication()
	if err != nil {
		log.NewError(operationID, "GetRecvFriendApplication failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "list", svrList, onServer, onLocal)

	aInBNot, bInANot, sameA, sameB := common.CheckFriendRequestDiff(onServer, onLocal)
	log.Debug(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := f.db.InsertFriendRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertFriendRequest failed ", err.Error())
			continue
		}
		callbackData := sdk.FriendApplicationAddedCallback(*onServer[index])
		f.friendListener.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
	}
	for _, index := range sameA {
		err := f.db.UpdateFriendRequest(onServer[index])
		if err != nil {
			if !strings.Contains(err.Error(), "RowsAffected == 0") {
				log.NewError(operationID, "UpdateFriendRequest failed ", err.Error(), *onServer[index])
				continue
			}
			if onServer[index].HandleResult == -1 {
				callbackData := sdk.FriendApplicationRejectCallback(*onServer[index])
				f.friendListener.OnFriendApplicationRejected(utils.StructToJsonString(callbackData))

			} else if onServer[index].HandleResult == 1 {
				callbackData := sdk.FriendApplicationAcceptCallback(*onServer[index])
				f.friendListener.OnFriendApplicationAccepted(utils.StructToJsonString(callbackData))
			}
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriendRequestBothUserID(onLocal[index].FromUserID, onLocal[index].ToUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriendRequestBothUserID failed ", err.Error(), onLocal[index].FromUserID, onLocal[index].ToUserID)
			continue
		}
		callbackData := sdk.FriendApplicationAcceptCallback(*onLocal[index])
		f.friendListener.OnFriendApplicationDeleted(utils.StructToJsonString(callbackData))
	}
}

func (f *Friend) SyncFriendList(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := f.getServerFriendList(operationID)
	if err != nil {
		log.NewError(operationID, "getServerFriendList failed ", err.Error())
		return
	}
	friendsInfoOnServer := common.TransferToLocalFriend(svrList)
	friendsInfoOnLocal, err := f.db.GetAllFriendList()
	if err != nil {
		log.NewError(operationID, "_getAllFriendList failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "list ", svrList, friendsInfoOnServer, friendsInfoOnLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckFriendListDiff(friendsInfoOnServer, friendsInfoOnLocal)
	log.NewInfo(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := f.db.InsertFriend(friendsInfoOnServer[index])
		if err != nil {
			log.NewError(operationID, "_insertFriend failed ", err.Error())
			continue
		}
		callbackData := sdk.FriendAddedCallback(*friendsInfoOnServer[index])
		f.friendListener.OnFriendAdded(utils.StructToJsonString(callbackData))
	}
	for _, index := range sameA {
		err := f.db.UpdateFriend(friendsInfoOnServer[index])
		if err != nil {
			if !strings.Contains(err.Error(), "RowsAffected == 0") {
				log.NewError(operationID, "UpdateFriendRequest failed ", err.Error(), *friendsInfoOnServer[index])
				continue
			}
			callbackData := sdk.FriendInfoChangedCallback(*friendsInfoOnLocal[index])
			f.friendListener.OnFriendInfoChanged(utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriend(friendsInfoOnLocal[index].FriendUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriend failed ", err.Error())
			continue
		}
		callbackData := sdk.FriendDeletedCallback(*friendsInfoOnLocal[index])
		f.friendListener.OnFriendDeleted(utils.StructToJsonString(callbackData))
	}
}

func (f *Friend) SyncBlackList(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := f.getServerBlackList(operationID)
	if err != nil {
		log.NewError(operationID, "getServerBlackList failed ", err.Error())
		return
	}
	blackListOnServer := common.TransferToLocalBlack(svrList, f.loginUserID)
	blackListOnLocal, err := f.db.GetBlackList()
	if err != nil {
		log.NewError(operationID, "_getBlackList failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "list ", svrList, blackListOnServer, blackListOnLocal)
	aInBNot, bInANot, sameA, _ := common.CheckBlackListDiff(blackListOnServer, blackListOnLocal)
	for _, index := range aInBNot {
		err := f.db.InsertBlack(blackListOnServer[index])
		if err != nil {
			log.NewError(operationID, "_insertFriend failed ", err.Error())
			continue
		}
		callbackData := sdk.BlackAddCallback(*blackListOnServer[index])
		f.friendListener.OnBlackAdded(utils.StructToJsonString(callbackData))
	}
	for _, index := range sameA {
		err := f.db.UpdateBlack(blackListOnServer[index])
		if err != nil {
			log.NewError(operationID, "_updateFriend failed ", err.Error())
			continue
		}
		//todo : add black info update callback
		log.Info(operationID, "black info update, do nothing ", blackListOnServer[index])
	}
	for _, index := range bInANot {
		err := f.db.DeleteBlack(blackListOnLocal[index].BlockUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriend failed ", err.Error())
			continue
		}
		callbackData := sdk.BlackAddCallback(*blackListOnLocal[index])
		f.friendListener.OnBlackDeleted(utils.StructToJsonString(callbackData))
	}
}

func (f *Friend) DoNotification(msg *api.MsgData) {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg)
	if f.friendListener == nil {
		log.Error(operationID, "f.friendListener == nil")
		return
	}

	go func() {
		switch msg.ContentType {
		case constant.FriendApplicationNotification:
			f.friendApplicationNotification(msg, operationID)
		case constant.FriendApplicationApprovedNotification:
			f.friendApplicationApprovedNotification(msg, operationID)
		case constant.FriendApplicationRejectedNotification:
			f.friendApplicationRejectedNotification(msg, operationID)
		case constant.FriendAddedNotification:
			f.friendAddedNotification(msg, operationID)
		case constant.FriendDeletedNotification:
			f.friendDeletedNotification(msg, operationID)
		case constant.FriendRemarkSetNotification:
			f.friendInfoChangedNotification(msg, operationID)
		case constant.UserInfoUpdatedNotification:
			f.friendInfoChangedNotification(msg, operationID)
		case constant.BlackAddedNotification:
			f.blackAddedNotification(msg, operationID)
		case constant.BlackDeletedNotification:
			f.blackDeletedNotification(msg, operationID)
		default:
			log.Error(operationID, "type failed ", msg.ClientMsgID, msg.ServerMsgID, msg.ContentType)
		}
	}()
}

func (f *Friend) blackDeletedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.BlackDeletedTips{FromToUserID: &api.FromToUserID{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg.Content)
		return
	}
	if detail.FromToUserID.FromUserID == f.loginUserID {
		f.SyncBlackList(operationID)
	}
}

func (f *Friend) blackAddedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.BlackAddedTips{FromToUserID: &api.FromToUserID{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg.Content)
		return
	}
	if detail.FromToUserID.FromUserID == f.loginUserID {
		f.SyncBlackList(operationID)
	}
}

func (f *Friend) friendInfoChangedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	var detail api.UserInfoUpdatedTips
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg.Content)
		return
	}
	if detail.UserID != f.loginUserID {
		f.SyncFriendList(operationID)
	}
}

func (f *Friend) friendDeletedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	f.SyncFriendList(operationID)
}

func (f *Friend) friendAddedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	f.SyncFriendList(operationID)
}

func (f *Friend) friendApplicationNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.FriendApplicationTips{FromToUserID: &api.FromToUserID{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg.Content)
		return
	}
	if detail.FromToUserID.FromUserID == f.loginUserID {
		f.SyncSelfFriendApplication(operationID)
	} else {
		f.SyncFriendApplication(operationID)
	}
}

func (f *Friend) friendApplicationRejectedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.FriendApplicationRejectedTips{FromToUserID: &api.FromToUserID{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg.Content)
		return
	}
	if f.loginUserID == detail.FromToUserID.FromUserID {
		f.SyncFriendApplication(operationID)
	} else {
		f.SyncSelfFriendApplication(operationID)
	}
}

func (f *Friend) friendApplicationApprovedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.FriendApplicationApprovedTips{FromToUserID: &api.FromToUserID{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg.Content)
		return
	}

	f.SyncFriendList(operationID)
	if f.loginUserID == detail.FromToUserID.FromUserID {
		f.SyncFriendApplication(operationID)
	} else {
		f.SyncSelfFriendApplication(operationID)
	}
}

//func (f *Friend) getLocalFriendList() ([]open_im_sdk.friendInfo, error) {
//	//Take out the friend list and judge whether it is in the blacklist again to prevent nested locks
//	localFriendList, err := u.getLocalFriendList22()
//	if err != nil {
//		return nil, err
//	}
//	for index, v := range localFriendList {
//		//check friend is in blacklist
//		blackUser, err := u.getBlackUsInfoByUid(v.UID)
//		if err != nil {
//			utils.sdkLog(err.Error())
//		}
//		if blackUser.Uid != "" {
//			localFriendList[index].IsInBlackList = 1
//		}
//	}
//	return localFriendList, nil
//}

//func (f *Friend) getServerSelfApplication() ([]api.applyUserInfo, error) {
//	resp, err := network.Post2Api(constant.GetSelfApplicationListRouter, open_im_sdk.paramsCommonReq{OperationID: utils.operationIDGenerator()}, u.token)
//	if err != nil {
//		return nil, err
//	}
//	var vgetFriendApplyListResp open_im_sdk.getFriendApplyListResp
//	err = json.Unmarshal(resp, &vgetFriendApplyListResp)
//	if err != nil {
//		utils.sdkLog("unmarshal failed, ", err.Error())
//		return nil, err
//	}
//	if vgetFriendApplyListResp.ErrCode != 0 {
//		utils.sdkLog("errcode: ", vgetFriendApplyListResp.ErrCode, "errmsg: ", vgetFriendApplyListResp.ErrMsg)
//		return nil, err
//	}
//	return vgetFriendApplyListResp.Data, nil
//}
