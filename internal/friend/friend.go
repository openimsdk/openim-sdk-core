package friend

import (
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/mapstructure"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"strings"
)

type Friend struct {
	friendListener OnFriendshipListener
	loginUserID    string
	db             *db.DataBase
	p              *ws.PostApi
}

func NewFriend(loginUserID string, db *db.DataBase, p *ws.PostApi) *Friend {
	return &Friend{loginUserID: loginUserID, db: db, p: p}
}

func (f *Friend) SetListener(listener OnFriendshipListener) {
	f.friendListener = listener
}

func (f *Friend) getDesignatedFriendsInfo(callback common.Base, friendUserIDList sdk_params_callback.GetDesignatedFriendsInfoParams, operationID string) sdk_params_callback.GetDesignatedFriendsInfoCallback {
	log.NewInfo(operationID, utils.GetSelfFuncName(), friendUserIDList)
	blackList, err := f.db.GetBlackInfoList(friendUserIDList)
	common.CheckErr(callback, err, operationID)
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
	common.CheckErr(callback, err, operationID)
	log.NewInfo(operationID, "_getFriendInfoList ", pureFriendUserIDList, localFriendList)
	return localFriendList
}

func (f *Friend) addFriend(callback common.Base, addFriendParams sdk_params_callback.AddFriendParams, operationID string) *server_api_params.CommDataResp {
	log.NewInfo(operationID, "addFriend args: ", addFriendParams)
	apiReq := server_api_params.AddFriendReq{}
	apiReq.ToUserID = addFriendParams.ToUserID
	apiReq.FromUserID = f.loginUserID
	apiReq.ReqMsg = addFriendParams.ReqMsg
	apiReq.OperationID = operationID
	return f.p.PostFatalCallback(callback, constant.AddFriendRouter, apiReq, operationID)
}

func (f *Friend) getRecvFriendApplicationList(callback common.Base, operationID string) sdk_params_callback.GetRecvFriendApplicationListCallback {
	log.NewInfo(operationID, "getRecvFriendApplicationList args: ")
	friendApplicationList, err := f.db.GetRecvFriendApplication()
	common.CheckErr(callback, err, operationID)
	return friendApplicationList
}

func (f *Friend) getSendFriendApplicationList(callback common.Base, operationID string) sdk_params_callback.GetSendFriendApplicationListCallback {
	return nil
}

func (f *Friend) processFriendApplication(callback common.Base, params sdk_params_callback.ProcessFriendApplicationParams, handleResult int32, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.AddFriendResponseReq{}
	apiReq.FromUserID = f.loginUserID
	apiReq.ToUserID = params.ToUserID
	apiReq.Flag = handleResult
	apiReq.OperationID = operationID
	apiReq.HandleMsg = params.HandleMsg
	result := f.p.PostFatalCallback(callback, constant.AddFriendResponse, apiReq, operationID)
	f.SyncFriendApplication()
	return result
}

func (f *Friend) checkFriend(callback common.Base, userIDList sdk_params_callback.CheckFriendParams, operationID string) sdk_params_callback.CheckFriendCallback {
	friendList, err := f.db.GetFriendInfoList(userIDList)
	common.CheckErr(callback, err, operationID)
	blackList, err := f.db.GetBlackInfoList(userIDList)
	common.CheckErr(callback, err, operationID)
	var checkFriendCallback sdk_params_callback.CheckFriendCallback
	for _, v := range userIDList {
		var r server_api_params.UserIDResult
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
	return checkFriendCallback
}

func (f *Friend) deleteFriend(FriendUserID string, callback common.Base, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.DeleteFriendReq{}
	apiReq.ToUserID = FriendUserID
	apiReq.FromUserID = f.loginUserID
	apiReq.OperationID = operationID
	result := f.p.PostFatalCallback(callback, constant.DeleteFriendRouter, apiReq, operationID)
	f.SyncFriendList()
	return result
}

func (f *Friend) setFriendRemark(params sdk_params_callback.SetFriendRemarkParams, callback common.Base, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.SetFriendRemarkReq{}
	apiReq.OperationID = operationID
	apiReq.ToUserID = params.ToUserID
	apiReq.FromUserID = f.loginUserID
	result := f.p.PostFatalCallback(callback, constant.SetFriendComment, apiReq, operationID)
	f.SyncFriendList()
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

func (f *Friend) getServerFriendList(operationID string) ([]*server_api_params.FriendInfo, error) {
	apiReq := server_api_params.GetFriendListReq{OperationID: operationID, FromUserID: f.loginUserID}
	resp, err := f.p.PostReturn(constant.GetFriendListRouter, apiReq)
	//	commData, err := common.CheckErrAndRespReturn(err, resp)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}

	realData := server_api_params.GetFriendListResp{}
	err = mapstructure.Decode(resp.Data, &realData.FriendInfoList)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}

	return realData.FriendInfoList, nil
}

func (f *Friend) getServerBlackList(operationID string) ([]*server_api_params.PublicUserInfo, error) {
	apiReq := server_api_params.GetBlackListReq{OperationID: operationID, FromUserID: f.loginUserID}
	commData, err := f.p.PostReturn(constant.GetBlackListRouter, apiReq)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	realData := server_api_params.GetBlackListResp{}
	mapstructure.Decode(commData.Data, &realData.BlackUserInfoList)
	return realData.BlackUserInfoList, nil
}

func (f *Friend) getServerFriendApplication(operationID string) ([]*server_api_params.FriendRequest, error) {
	apiReq := server_api_params.GetFriendApplyListReq{OperationID: operationID, FromUserID: f.loginUserID}
	resp, err := f.p.PostReturn(constant.GetFriendApplicationListRouter, apiReq)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	realData := server_api_params.GetFriendApplyListResp{}
	mapstructure.Decode(resp.Data, &realData.FriendRequestList)
	return realData.FriendRequestList, nil
}

//func (f *Friend) getServerSelfApplication() ([]server_api_params.applyUserInfo, error) {
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
func (f *Friend) addBlack(callback common.Base, blackUid, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.AddBlacklistReq{}
	apiReq.ToUserID = blackUid
	apiReq.FromUserID = f.loginUserID
	apiReq.OperationID = operationID
	result := f.p.PostFatalCallback(callback, constant.AddBlackListRouter, apiReq, operationID)
	f.SyncBlackList()
	return result
}

func (f *Friend) removeBlack(callback common.Base, deleteUid, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.RemoveBlackListReq{}
	apiReq.ToUserID = deleteUid
	apiReq.FromUserID = f.loginUserID
	apiReq.OperationID = operationID
	result := f.p.PostFatalCallback(callback, constant.RemoveBlackListRouter, apiReq, operationID)
	f.SyncBlackList()
	return result
}

func (f *Friend) SyncSelfFriendApplication() {

}

func (f *Friend) SyncFriendApplication() {
	operationID := utils.OperationIDGenerator()
	svrList, err := f.getServerFriendApplication(operationID)
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
		callbackData := sdk_params_callback.FriendApplicationListAddedCallback(*onServer[index])
		f.friendListener.OnFriendApplicationListAdded(utils.StructToJsonString(callbackData))
	}
	for _, index := range sameA {
		err := f.db.UpdateFriendRequest(onServer[index])
		if err != nil {
			if !strings.Contains(err.Error(), "RowsAffected == 0") {
				log.NewError(operationID, "UpdateFriendRequest failed ", err.Error(), *onServer[index])
				continue
			}
			if onServer[index].HandleResult == -1 {
				callbackData := sdk_params_callback.FriendApplicationListRejectCallback(*onServer[index])
				f.friendListener.OnFriendApplicationListReject(utils.StructToJsonString(callbackData))

			} else if onServer[index].HandleResult == -1 {
				callbackData := sdk_params_callback.FriendApplicationListAcceptCallback(*onServer[index])
				f.friendListener.OnFriendApplicationListAccept(utils.StructToJsonString(callbackData))
			}
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriendRequestBothUserID(onLocal[index].FromUserID, onLocal[index].ToUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriendRequestBothUserID failed ", err.Error(), onLocal[index].FromUserID, onLocal[index].ToUserID)
			continue
		}
		callbackData := sdk_params_callback.FriendApplicationListAcceptCallback(*onLocal[index])
		f.friendListener.OnFriendApplicationListDeleted(utils.StructToJsonString(callbackData))
	}
}

func (f *Friend) SyncFriendList() {
	operationID := utils.OperationIDGenerator()
	svrList, err := f.getServerFriendList(operationID)
	if err != nil {
		log.NewError(operationID, "getServerFriendList failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "svrList", svrList)
	friendsInfoOnServer := common.TransferToLocalFriend(svrList)
	friendsInfoOnLocal, err := f.db.GetAllFriendList()
	if err != nil {
		log.NewError(operationID, "_getAllFriendList failed ", err.Error())
		return
	}

	log.NewInfo(operationID, "friendsInfoOnServer", friendsInfoOnServer)
	log.NewInfo(operationID, "friendsInfoOnLocal", friendsInfoOnLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckFriendListDiff(friendsInfoOnServer, friendsInfoOnLocal)
	log.NewInfo(operationID, "checkFriendListDiff", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := f.db.InsertFriend(friendsInfoOnServer[index])
		if err != nil {
			log.NewError(operationID, "_insertFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := f.db.UpdateFriend(friendsInfoOnServer[index])
		if err != nil {
			log.NewError(operationID, "_updateFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriend(friendsInfoOnLocal[index].FriendUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriend failed ", err.Error())
			continue
		}
	}
}

func (f *Friend) SyncBlackList() {
	operationID := utils.OperationIDGenerator()
	svrList, err := f.getServerBlackList(operationID)
	if err != nil {
		log.NewError(operationID, "getServerBlackList failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "svrList", svrList)
	blackListOnServer := common.TransferToLocalBlack(svrList, f.loginUserID)
	blackListOnLocal, err := f.db.GetBlackList()
	if err != nil {
		log.NewError(operationID, "_getBlackList failed ", err.Error())
		return
	}

	log.NewInfo(operationID, "blackListOnServer", blackListOnServer)
	log.NewInfo(operationID, "blackListOnlocal", blackListOnLocal)
	aInBNot, bInANot, sameA, _ := common.CheckBlackListDiff(blackListOnServer, blackListOnLocal)
	for _, index := range aInBNot {
		err := f.db.InsertBlack(blackListOnServer[index])
		if err != nil {
			log.NewError(operationID, "_insertFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := f.db.UpdateBlack(blackListOnServer[index])
		if err != nil {
			log.NewError(operationID, "_updateFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteBlack(blackListOnLocal[index].BlockUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriend failed ", err.Error())
			continue
		}
	}

}

func (u *Friend) DoFriendNotification(msg *server_api_params.MsgData) {
	if u.friendListener == nil {
		return
	}

	if msg.SendID == u.loginUserID && msg.SenderPlatformID == sdk_struct.SvrConf.Platform {
		return
	}

	go func() {
		switch msg.ContentType {
		case constant.FriendApplicationProcessedNotification:
			u.friendApplicationProcessedNotification(msg)
		case constant.FriendApplicationAddedNotification:
			u.friendApplicationAddedNotification(msg)
		case constant.FriendAddedNotification:
			u.friendAddedNotification(msg)
		case constant.FriendDeletedNotification:
			u.friendDeletedNotification(msg)
		case constant.FriendInfoChangedNotification:
			u.friendInfoChangedNotification(msg)
		case constant.BlackAddedNotification:
			u.blackAddedNotification(msg)
		case constant.BlackDeletedNotification:
			u.blackDeletedNotification(msg)
		default:
			log.Error("", "type failed ", msg.ClientMsgID, msg.ServerMsgID, msg.ContentType)
		}
	}()
}

func (u *Friend) blackDeletedNotification(msg *server_api_params.MsgData) {
	u.SyncBlackList()
}

func (u *Friend) blackAddedNotification(msg *server_api_params.MsgData) {
	u.SyncBlackList()
}

func (u *Friend) friendInfoChangedNotification(msg *server_api_params.MsgData) {
	u.SyncFriendList()
}

func (u *Friend) friendDeletedNotification(msg *server_api_params.MsgData) {
	u.SyncFriendList()
}

func (u *Friend) friendAddedNotification(msg *server_api_params.MsgData) {
	u.SyncFriendList()
}

func (u *Friend) friendApplicationAddedNotification(msg *server_api_params.MsgData) {
	u.SyncFriendApplication()
	u.SyncSelfFriendApplication()
}

func (u *Friend) friendApplicationProcessedNotification(msg *server_api_params.MsgData) {
	var tips server_api_params.TipsComm
	proto.Unmarshal(msg.Content, &tips)

	var detail server_api_params.FriendApplicationProcessedTips
	proto.Unmarshal(tips.Detail, &detail)

	u.SyncFriendList()
	u.SyncFriendApplication()
	u.SyncSelfFriendApplication()
}
