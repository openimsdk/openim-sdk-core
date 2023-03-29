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
	"errors"
	comm "open_im_sdk/internal/common"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/internal/user"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

type Friend struct {
	friendListener open_im_sdk_callback.OnFriendshipListener
	loginUserID    string
	db             db_interface.DataBase
	user           *user.User
	p              *ws.PostApi
	loginTime      int64
	conversationCh chan common.Cmd2Value

	listenerForService open_im_sdk_callback.OnListenerForService
}

func (f *Friend) LoginTime() int64 {
	return f.loginTime
}

func (f *Friend) SetLoginTime(loginTime int64) {
	f.loginTime = loginTime
}

func (f *Friend) Db() db_interface.DataBase {
	return f.db
}

func NewFriend(loginUserID string, db db_interface.DataBase, user *user.User, p *ws.PostApi, conversationCh chan common.Cmd2Value) *Friend {
	return &Friend{loginUserID: loginUserID, db: db, user: user, p: p, conversationCh: conversationCh}
}

func (f *Friend) SetListener(listener open_im_sdk_callback.OnFriendshipListener) {
	f.friendListener = listener
}

func (f *Friend) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	f.listenerForService = listener
}

func (f *Friend) getDesignatedFriendsInfo(callback open_im_sdk_callback.Base, friendUserIDList sdk.GetDesignatedFriendsInfoParams, operationID string) sdk.GetDesignatedFriendsInfoCallback {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", friendUserIDList)

	localFriendList, err := f.db.GetFriendInfoList(friendUserIDList)
	common.CheckDBErrCallback(callback, err, operationID)

	blackList, err := f.db.GetBlackInfoList(friendUserIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	for _, v := range blackList {
		log.Info(operationID, "GetBlackInfoList ", *v)
	}

	r := common.MergeFriendBlackResult(localFriendList, blackList)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", r)
	return r
}

func (f *Friend) GetUserNameAndFaceUrlByUid(friendUserID, operationID string) (faceUrl, name string, err error, isFromSvr bool) {
	isFromSvr = false
	friendInfo, err := f.db.GetFriendInfoByFriendUserID(friendUserID)
	if err == nil {
		if friendInfo.Remark != "" {
			return friendInfo.FaceURL, friendInfo.Remark, nil, isFromSvr
		} else {
			return friendInfo.FaceURL, friendInfo.Nickname, nil, isFromSvr
		}
	} else {
		if operationID == "" {
			operationID = utils.OperationIDGenerator()
		}
		userInfos, err := f.user.GetUsersInfoFromSvrNoCallback([]string{friendUserID}, operationID)
		if err != nil {
			return "", "", err, isFromSvr
		}
		for _, v := range userInfos {
			isFromSvr = true
			return v.FaceURL, v.Nickname, nil, isFromSvr
		}
		log.Info(operationID, "GetUsersInfoFromSvr ", friendUserID)
	}
	return "", "", errors.New("getUserNameAndFaceUrlByUid err"), isFromSvr
}

func (f *Friend) GetDesignatedFriendListInfo(callback open_im_sdk_callback.Base, friendUserIDList []string, operationID string) []*model_struct.LocalFriend {
	friendList, err := f.db.GetFriendInfoList(friendUserIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	return friendList
}

func (f *Friend) GetDesignatedBlackListInfo(callback open_im_sdk_callback.Base, blackIDList []string, operationID string) []*model_struct.LocalBlack {
	blackList, err := f.db.GetBlackInfoList(blackIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	return blackList
}

func (f *Friend) addFriend(callback open_im_sdk_callback.Base, userIDReqMsg sdk.AddFriendParams, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", userIDReqMsg)
	apiReq := api.AddFriendReq{}
	apiReq.ToUserID = userIDReqMsg.ToUserID
	apiReq.FromUserID = f.loginUserID
	apiReq.ReqMsg = userIDReqMsg.ReqMsg
	apiReq.OperationID = operationID
	f.p.PostFatalCallbackPenetrate(callback, constant.AddFriendRouter, apiReq, nil, operationID)
	f.SyncFriendApplication(operationID)
}

func (f *Friend) getRecvFriendApplicationList(callback open_im_sdk_callback.Base, operationID string) sdk.GetRecvFriendApplicationListCallback {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	friendApplicationList, err := f.db.GetRecvFriendApplication()
	common.CheckDBErrCallback(callback, err, operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", friendApplicationList)
	return friendApplicationList
}

func (f *Friend) getSendFriendApplicationList(callback open_im_sdk_callback.Base, operationID string) sdk.GetSendFriendApplicationListCallback {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	friendApplicationList, err := f.db.GetSendFriendApplication()
	common.CheckDBErrCallback(callback, err, operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", friendApplicationList)
	return friendApplicationList
}

func (f *Friend) processFriendApplication(callback open_im_sdk_callback.Base, userIDHandleMsg sdk.ProcessFriendApplicationParams, handleResult int32, operationID string) {
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

func (f *Friend) checkFriend(callback open_im_sdk_callback.Base, friendUserIDList sdk.CheckFriendParams, operationID string) sdk.CheckFriendCallback {
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

func (f *Friend) deleteFriend(friendUserID sdk.DeleteFriendParams, callback open_im_sdk_callback.Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", friendUserID)
	apiReq := api.DeleteFriendReq{}
	apiReq.ToUserID = string(friendUserID)
	apiReq.FromUserID = f.loginUserID
	apiReq.OperationID = operationID
	f.p.PostFatalCallback(callback, constant.DeleteFriendRouter, apiReq, nil, operationID)
	f.SyncFriendList(operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ")
}

func (f *Friend) getFriendList(callback open_im_sdk_callback.Base, operationID string) sdk.GetFriendListCallback {
	localFriendList, err := f.db.GetAllFriendList()
	common.CheckDBErrCallback(callback, err, operationID)
	localBlackList, err := f.db.GetBlackListDB()
	common.CheckDBErrCallback(callback, err, operationID)
	return common.MergeFriendBlackResult(localFriendList, localBlackList)
}
func (f *Friend) searchFriends(callback open_im_sdk_callback.Base, param sdk.SearchFriendsParam, operationID string) sdk.SearchFriendsCallback {
	if len(param.KeywordList) == 0 || (!param.IsSearchNickname && !param.IsSearchUserID && !param.IsSearchRemark) {
		common.CheckAnyErrCallback(callback, 201, errors.New("keyword is null or search field all false"), operationID)
	}
	localFriendList, err := f.db.SearchFriendList(param.KeywordList[0], param.IsSearchUserID, param.IsSearchNickname, param.IsSearchRemark)
	common.CheckDBErrCallback(callback, err, operationID)
	localBlackList, err := f.db.GetBlackListDB()
	common.CheckDBErrCallback(callback, err, operationID)
	return mergeFriendBlackSearchResult(localFriendList, localBlackList)
}
func mergeFriendBlackSearchResult(base []*model_struct.LocalFriend, add []*model_struct.LocalBlack) (result []*sdk.SearchFriendItem) {
	blackUserIDList := func(bl []*model_struct.LocalBlack) (result []string) {
		for _, v := range bl {
			result = append(result, v.BlockUserID)
		}
		return result
	}(add)
	for _, v := range base {
		node := sdk.SearchFriendItem{}
		node.OwnerUserID = v.OwnerUserID
		node.FriendUserID = v.FriendUserID
		node.Remark = v.Remark
		node.CreateTime = v.CreateTime
		node.AddSource = v.AddSource
		node.OperatorUserID = v.OperatorUserID
		node.Nickname = v.Nickname
		node.FaceURL = v.FaceURL
		node.Gender = v.Gender
		node.PhoneNumber = v.PhoneNumber
		node.Birth = v.Birth
		node.Email = v.Email
		node.Ex = v.Ex
		node.AttachedInfo = v.AttachedInfo
		if !utils.IsContain(v.FriendUserID, blackUserIDList) {
			node.Relationship = constant.FriendRelationship
		}
		result = append(result, &node)
	}
	return result
}
func (f *Friend) getBlackList(callback open_im_sdk_callback.Base, operationID string) sdk.GetBlackListCallback {
	localBlackList, err := f.db.GetBlackListDB()
	common.CheckDBErrCallback(callback, err, operationID)

	localFriendList, err := f.db.GetAllFriendList()
	common.CheckDBErrCallback(callback, err, operationID)

	return common.MergeBlackFriendResult(localBlackList, localFriendList)
}

func (f *Friend) setFriendRemark(userIDRemark sdk.SetFriendRemarkParams, callback open_im_sdk_callback.Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", userIDRemark)
	apiReq := api.SetFriendRemarkReq{}
	apiReq.OperationID = operationID
	apiReq.ToUserID = userIDRemark.ToUserID
	apiReq.FromUserID = f.loginUserID
	apiReq.Remark = userIDRemark.Remark
	f.p.PostFatalCallback(callback, constant.SetFriendRemark, apiReq, nil, operationID)
	f.SyncFriendList(operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ")
}

func (f *Friend) getServerFriendList(operationID string) ([]*api.FriendInfo, error) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	apiReq := api.GetFriendListReq{OperationID: operationID, FromUserID: f.loginUserID}
	realData := api.GetFriendListResp{}
	err := f.p.PostReturn(constant.GetFriendListRouter, apiReq, &realData.FriendInfoList)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", realData.FriendInfoList)
	return realData.FriendInfoList, nil
}

func (f *Friend) getServerBlackList(operationID string) ([]*api.PublicUserInfo, error) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	apiReq := api.GetBlackListReq{OperationID: operationID, FromUserID: f.loginUserID}
	realData := api.GetBlackListResp{}
	err := f.p.PostReturn(constant.GetBlackListRouter, apiReq, &realData.BlackUserInfoList)
	if err != nil {
		return nil, utils.Wrap(err, operationID)
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", realData.BlackUserInfoList)
	return realData.BlackUserInfoList, nil
}

// recv
func (f *Friend) getFriendApplicationFromServer(operationID string) ([]*api.FriendRequest, error) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	apiReq := api.GetFriendApplyListReq{OperationID: operationID, FromUserID: f.loginUserID}
	realData := api.GetFriendApplyListResp{}
	err := f.p.PostReturn(constant.GetFriendApplicationListRouter, apiReq, &realData.FriendRequestList)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", realData.FriendRequestList)
	return realData.FriendRequestList, nil
}

// send
func (f *Friend) getSelfFriendApplicationFromServer(operationID string) ([]*api.FriendRequest, error) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	apiReq := api.GetSelfFriendApplyListReq{OperationID: operationID, FromUserID: f.loginUserID}
	realData := api.GetSelfFriendApplyListResp{}
	err := f.p.PostReturn(constant.GetSelfFriendApplicationListRouter, apiReq, &realData.FriendRequestList)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", realData.FriendRequestList)
	return realData.FriendRequestList, nil
}

func (f *Friend) addBlack(callback open_im_sdk_callback.Base, blackUserID sdk.AddBlackParams, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", blackUserID)
	apiReq := api.AddBlacklistReq{}
	apiReq.ToUserID = string(blackUserID)
	apiReq.FromUserID = f.loginUserID
	apiReq.OperationID = operationID
	f.p.PostFatalCallback(callback, constant.AddBlackRouter, apiReq, nil, operationID)
	f.SyncBlackList(operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ")
}

func (f *Friend) removeBlack(callback open_im_sdk_callback.Base, blackUserID sdk.RemoveBlackParams, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", blackUserID)
	apiReq := api.RemoveBlackListReq{}
	apiReq.ToUserID = string(blackUserID)
	apiReq.FromUserID = f.loginUserID
	apiReq.OperationID = operationID
	f.p.PostFatalCallback(callback, constant.RemoveBlackRouter, apiReq, nil, operationID)
	f.SyncBlackList(operationID)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ")
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
		if f.friendListener != nil {
			f.friendListener.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnFriendApplicationAdded", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range sameA {
		err := f.db.UpdateFriendRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateFriendRequest failed ", err.Error(), *onServer[index])
			continue
		} else {
			if onServer[index].HandleResult == constant.FriendResponseRefuse {
				callbackData := sdk.FriendApplicationRejectCallback(*onServer[index])
				if f.friendListener != nil {

					f.friendListener.OnFriendApplicationRejected(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnFriendApplicationRejected", utils.StructToJsonString(callbackData))
				}

			} else if onServer[index].HandleResult == constant.FriendResponseAgree {
				callbackData := sdk.FriendApplicationAcceptCallback(*onServer[index])
				if f.friendListener != nil {
					f.friendListener.OnFriendApplicationAccepted(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnFriendApplicationAccepted", utils.StructToJsonString(callbackData))
				}
				if f.listenerForService != nil {
					f.listenerForService.OnFriendApplicationAccepted(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnFriendApplicationAccepted", utils.StructToJsonString(callbackData))
				}
			} else {
				callbackData := sdk.FriendApplicationAddedCallback(*onServer[index])
				if f.friendListener != nil {
					f.friendListener.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnFriendApplicationAdded", utils.StructToJsonString(callbackData))
				}
			}
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriendRequestBothUserID(onLocal[index].FromUserID, onLocal[index].ToUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriendRequestBothUserID failed ", err.Error(), onLocal[index].FromUserID, onLocal[index].ToUserID)
			continue
		}
		callbackData := sdk.FriendApplicationDeletedCallback(*onLocal[index])
		if f.friendListener != nil {

			f.friendListener.OnFriendApplicationDeleted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnFriendApplicationDeleted", utils.StructToJsonString(callbackData))
		}
	}
}

// recv
func (f *Friend) SyncFriendApplication(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := f.getFriendApplicationFromServer(operationID)
	if err != nil {
		log.NewError(operationID, "getFriendApplicationFromServer failed ", err.Error())
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
		//f.friendListener.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
		if f.friendListener != nil {
			f.friendListener.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnReceiveFriendApplicationAdded", utils.StructToJsonString(callbackData))
		}
		if f.listenerForService != nil {
			f.listenerForService.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnReceiveFriendApplicationAdded", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range sameA {
		err := f.db.UpdateFriendRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateFriendRequest failed ", err.Error(), *onServer[index])
			continue
		} else {
			if onServer[index].HandleResult == constant.FriendResponseRefuse {
				callbackData := sdk.FriendApplicationRejectCallback(*onServer[index])
				if f.friendListener != nil {

					f.friendListener.OnFriendApplicationRejected(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnFriendApplicationRejected", utils.StructToJsonString(callbackData))
				}
			} else if onServer[index].HandleResult == constant.FriendResponseAgree {
				callbackData := sdk.FriendApplicationAcceptCallback(*onServer[index])
				if f.friendListener != nil {

					f.friendListener.OnFriendApplicationAccepted(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnFriendApplicationAccepted", utils.StructToJsonString(callbackData))
				}
			} else {
				callbackData := sdk.FriendApplicationAddedCallback(*onServer[index])
				if f.friendListener != nil {
					f.friendListener.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnReceiveFriendApplicationAdded", utils.StructToJsonString(callbackData))
				}
				if f.listenerForService != nil {
					f.listenerForService.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnReceiveFriendApplicationAdded", utils.StructToJsonString(callbackData))
				}
			}
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriendRequestBothUserID(onLocal[index].FromUserID, onLocal[index].ToUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriendRequestBothUserID failed ", err.Error(), onLocal[index].FromUserID, onLocal[index].ToUserID)
			continue
		}
		callbackData := sdk.FriendApplicationDeletedCallback(*onLocal[index])
		if f.friendListener != nil {

			f.friendListener.OnFriendApplicationDeleted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnReceiveFriendApplicationDeleted", utils.StructToJsonString(callbackData))
		}
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
	for _, v := range friendsInfoOnServer {
		log.NewDebug(operationID, "friendsInfoOnServer ", *v)
	}
	aInBNot, bInANot, sameA, sameB := common.CheckFriendListDiff(friendsInfoOnServer, friendsInfoOnLocal)
	log.NewInfo(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := f.db.InsertFriend(friendsInfoOnServer[index])
		if err != nil {
			log.NewError(operationID, "_insertFriend failed ", err.Error())
			continue
		}
		callbackData := sdk.FriendAddedCallback(*friendsInfoOnServer[index])
		if f.friendListener != nil {

			f.friendListener.OnFriendAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnFriendAdded", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range sameA {
		callbackData := sdk.FriendInfoChangedCallback(*friendsInfoOnServer[index])
		localFriend, err := f.db.GetFriendInfoByFriendUserID(callbackData.FriendUserID)
		if err != nil {
			log.NewError(operationID, "GetFriendInfoByFriendUserID failed ", err.Error(), "userID", callbackData.FriendUserID)
			continue
		}
		err = f.db.UpdateFriend(friendsInfoOnServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateFriendRequest failed ", err.Error(), *friendsInfoOnServer[index])
			continue
		} else {
			callbackData := sdk.FriendInfoChangedCallback(*friendsInfoOnServer[index])
			if f.friendListener != nil {
				f.friendListener.OnFriendInfoChanged(utils.StructToJsonString(callbackData))
				if localFriend.Nickname == callbackData.Nickname && localFriend.FaceURL == callbackData.FaceURL && localFriend.Remark == callbackData.Remark {
					log.NewInfo(operationID, "OnFriendInfoChanged nickname faceURL unchanged", callbackData.FriendUserID, localFriend.Nickname, localFriend.FaceURL)
					continue
				}
				if callbackData.Remark != "" {
					callbackData.Nickname = callbackData.Remark
				}
				common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName, Args: common.SourceIDAndSessionType{SourceID: callbackData.FriendUserID, SessionType: constant.SingleChatType}}, f.conversationCh)
				common.TriggerCmdUpdateMessage(common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName, Args: common.UpdateMessageInfo{UserID: callbackData.FriendUserID, FaceURL: callbackData.FaceURL, Nickname: callbackData.Nickname}}, f.conversationCh)
				log.Info(operationID, "OnFriendInfoChanged", utils.StructToJsonString(callbackData))
			}
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriendDB(friendsInfoOnLocal[index].FriendUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriend failed ", err.Error())
			continue
		}
		callbackData := sdk.FriendDeletedCallback(*friendsInfoOnLocal[index])
		if f.friendListener != nil {

			f.friendListener.OnFriendDeleted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnFriendDeleted", utils.StructToJsonString(callbackData))
		}
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
	blackListOnLocal, err := f.db.GetBlackListDB()
	if err != nil {
		log.NewError(operationID, "_getBlackList failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "list ", svrList, blackListOnServer, blackListOnLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckBlackListDiff(blackListOnServer, blackListOnLocal)
	log.NewInfo(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := f.db.InsertBlack(blackListOnServer[index])
		if err != nil {
			log.NewError(operationID, "_insertFriend failed ", err.Error())
			continue
		}
		callbackData := sdk.BlackAddCallback(*blackListOnServer[index])
		if f.friendListener != nil {

			f.friendListener.OnBlackAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnBlackAdded", utils.StructToJsonString(callbackData))
		}
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
		callbackData := sdk.BlackDeletedCallback(*blackListOnLocal[index])
		if f.friendListener != nil {
			f.friendListener.OnBlackDeleted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnBlackDeleted", utils.StructToJsonString(callbackData))
		}
	}
}

func (f *Friend) DoNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value) {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg)
	if f.friendListener == nil {
		log.Error(operationID, "f.friendListener == nil")
		return
	}
	if msg.SendTime < f.loginTime || f.loginTime == 0 {
		log.Warn(operationID, "ignore notification ", msg.ClientMsgID, msg.ServerMsgID, msg.Seq, msg.ContentType)
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
			f.friendRemarkNotification(msg, conversationCh, operationID)
		case constant.FriendInfoUpdatedNotification:
			f.friendInfoChangedNotification(msg, conversationCh, operationID)
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

func (f *Friend) friendRemarkNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	var detail api.FriendInfoChangedTips
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg.Content)
		return
	}
	if detail.FromToUserID.FromUserID == f.loginUserID {
		f.SyncFriendList(operationID)
	}
}

func (f *Friend) friendInfoChangedNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	var detail api.UserInfoUpdatedTips
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg.Content)
		return
	}
	if detail.UserID != f.loginUserID {
		f.SyncFriendList(operationID)
	} else {
		log.Warn(operationID, "detail failed,  detail.UserID == f.loginUserID ", detail.UserID)
		f.SyncFriendList(operationID)
		//f.user.SyncLoginUserInfo(operationID)
		//go func() {
		//	loginUserInfo, err := f.db.GetLoginUser(f.loginUserID)
		//	if err == nil {
		//		//_ = f.db.UpdateMsgSenderFaceURLAndSenderNickname(detail.UserID, loginUserInfo.FaceURL, loginUserInfo.Nickname, constant.SingleChatType)
		//		_ = common.TriggerCmdUpdateMessage(common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName, Args: common.UpdateMessageInfo{UserID: detail.UserID, FaceURL: loginUserInfo.FaceURL, Nickname: loginUserInfo.Nickname}}, conversationCh)
		//
		//	}
		//}()
	}
}

func (f *Friend) friendDeletedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.FriendDeletedTips{FromToUserID: &api.FromToUserID{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg)
		return
	}
	if detail.FromToUserID.FromUserID == f.loginUserID {
		f.SyncFriendList(operationID)
		return
	}
}

func (f *Friend) friendAddedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.FriendAddedTips{Friend: &api.FriendInfo{}, OpUser: &api.PublicUserInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg)
		return
	}
	log.Info("detail: ", detail.Friend)
	f.SyncFriendList(operationID)
	if detail.Friend.OwnerUserID == f.loginUserID || detail.Friend.FriendUser.UserID == f.loginUserID {
		f.SyncFriendList(operationID)
		return
	}
}

func (f *Friend) friendApplicationNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.FriendApplicationTips{FromToUserID: &api.FromToUserID{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg)
		return
	}
	if detail.FromToUserID.FromUserID == f.loginUserID {
		log.Info(operationID, "SyncSelfFriendApplication ", detail.FromToUserID.FromUserID)
		f.SyncSelfFriendApplication(operationID)
		return
	}
	if detail.FromToUserID.ToUserID == f.loginUserID {
		log.Info(operationID, "SyncFriendApplication ", detail.FromToUserID.FromUserID, detail.FromToUserID.ToUserID)
		f.SyncFriendApplication(operationID)
		return
	}
	log.Error(operationID, "FromToUserID failed ", detail.FromToUserID)
}

func (f *Friend) friendApplicationRejectedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.FriendApplicationRejectedTips{FromToUserID: &api.FromToUserID{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg)
		return
	}
	if f.loginUserID == detail.FromToUserID.FromUserID {
		f.SyncFriendApplication(operationID)
		return
	}
	if f.loginUserID == detail.FromToUserID.ToUserID {
		f.SyncSelfFriendApplication(operationID)
		return
	}
	log.Error(operationID, "FromToUserID failed ", detail.FromToUserID)
}

func (f *Friend) friendApplicationApprovedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.FriendApplicationApprovedTips{FromToUserID: &api.FromToUserID{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg)
		return
	}

	//f.SyncFriendList(operationID)
	if f.loginUserID == detail.FromToUserID.FromUserID {
		f.SyncFriendApplication(operationID)
		return
	}
	if f.loginUserID == detail.FromToUserID.ToUserID {
		f.SyncSelfFriendApplication(operationID)
		return
	}
	log.Error(operationID, "FromToUserID failed ", detail.FromToUserID)
}
