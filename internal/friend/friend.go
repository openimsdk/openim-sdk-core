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
	"context"
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
	"open_im_sdk/pkg/syncer"
	"open_im_sdk/pkg/utils"
)

func NewFriend(loginUserID string, db db_interface.DataBase, user *user.User, p *ws.PostApi, conversationCh chan common.Cmd2Value) *Friend {
	f := &Friend{loginUserID: loginUserID, db: db, user: user, p: p, conversationCh: conversationCh}
	f.initSyncer()
	return f
}

type Friend struct {
	friendListener     open_im_sdk_callback.OnFriendshipListener
	loginUserID        string
	db                 db_interface.DataBase
	user               *user.User
	p                  *ws.PostApi
	friendSyncer       *syncer.Syncer[*model_struct.LocalFriend, [2]string]
	blockSyncer        *syncer.Syncer[*model_struct.LocalBlack, [2]string]
	requestRecvSyncer  *syncer.Syncer[*model_struct.LocalFriendRequest, string]
	requestSendSyncer  *syncer.Syncer[*model_struct.LocalFriendRequest, string]
	loginTime          int64
	conversationCh     chan common.Cmd2Value
	listenerForService open_im_sdk_callback.OnListenerForService
}

func (f *Friend) initSyncer() {
	f.friendSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalFriend) error {
		return f.db.InsertFriend(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalFriend) error {
		return f.db.DeleteFriendDB(ctx, value.FriendUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriend, local *model_struct.LocalFriend) error {
		return f.db.UpdateFriend(ctx, server)
	}, func(value *model_struct.LocalFriend) [2]string {
		return [...]string{value.OwnerUserID, value.FriendUserID}
	}, nil, nil)

	f.blockSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalBlack) error {
		return f.db.InsertBlack(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalBlack) error {
		return f.db.DeleteBlack(ctx, value.BlockUserID)
	}, func(ctx context.Context, server *model_struct.LocalBlack, local *model_struct.LocalBlack) error {
		return f.db.UpdateBlack(ctx, server)
	}, func(value *model_struct.LocalBlack) [2]string {
		return [...]string{value.OwnerUserID, value.BlockUserID}
	}, nil, nil)

	f.requestRecvSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.InsertFriendRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.DeleteFriendRequestBothUserID(ctx, value.FromUserID, value.ToUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriendRequest, local *model_struct.LocalFriendRequest) error {
		return f.db.UpdateFriendRequest(ctx, server)
	}, func(value *model_struct.LocalFriendRequest) string {
		return value.FromUserID
	}, nil, nil)

	f.requestSendSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.InsertFriendRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.DeleteFriendRequestBothUserID(ctx, value.FromUserID, value.ToUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriendRequest, local *model_struct.LocalFriendRequest) error {
		return f.db.UpdateFriendRequest(ctx, server)
	}, func(value *model_struct.LocalFriendRequest) string {
		return value.ToUserID
	}, nil, nil)

}

func (f *Friend) LoginTime() int64 {
	return f.loginTime
}

func (f *Friend) SetLoginTime(loginTime int64) {
	f.loginTime = loginTime
}

func (f *Friend) SetFriendListener(listener open_im_sdk_callback.OnFriendshipListener) {
	if listener == nil {
		return
	}
	f.friendListener = listener
}

func (f *Friend) Db() db_interface.DataBase {
	return f.db
}

func (f *Friend) SetListener(listener open_im_sdk_callback.OnFriendshipListener) {
	f.friendListener = listener
}

func (f *Friend) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	f.listenerForService = listener
}

//func (f *Friend) getDesignatedFriendsInfo(callback open_im_sdk_callback.Base, friendUserIDList sdk.GetDesignatedFriendsInfoParams, operationID string) sdk.GetDesignatedFriendsInfoCallback {
//	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", friendUserIDList)
//
//	localFriendList, err := f.db.GetFriendInfoList(friendUserIDList)
//	common.CheckDBErrCallback(callback, err, operationID)
//
//	blackList, err := f.db.GetBlackInfoList(friendUserIDList)
//	common.CheckDBErrCallback(callback, err, operationID)
//	for _, v := range blackList {
//		log.Info(operationID, "GetBlackInfoList ", *v)
//	}
//
//	r := common.MergeFriendBlackResult(localFriendList, blackList)
//	log.NewInfo(operationID, utils.GetSelfFuncName(), "return: ", r)
//	return r
//}

func (f *Friend) GetUserNameAndFaceUrlByUid(ctx context.Context, friendUserID string) (faceUrl, name string, err error, isFromSvr bool) {
	isFromSvr = false
	friendInfo, err := f.db.GetFriendInfoByFriendUserID(ctx, friendUserID)
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
