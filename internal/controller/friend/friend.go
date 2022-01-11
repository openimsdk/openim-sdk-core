package friend

import (
	"github.com/mitchellh/mapstructure"
	ws "open_im_sdk/internal/controller/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/network"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

type Friend struct {
	friendListener OnFriendshipListener
	token          string
	loginUserID    string
	db             *db.DataBase
	p              *ws.PostApi
}

func (f *Friend) Init(token, userID string, db *db.DataBase) {
	f.token = token
	f.loginUserID = userID
	f.db = db
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
	log.NewInfo(apiReq.OperationID, "post2Api ", constant.AddFriendRouter, apiReq)
	return f.p.PostFatalCallback(callback, constant.AddFriendRouter, apiReq, f.token)
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
	result := f.p.PostFatalCallback(callback, constant.AddFriendResponse, apiReq, f.token)
	f.syncFriendApplication()
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
	result := f.p.PostFatalCallback(callback, constant.DeleteFriendRouter, apiReq, f.token)
	f.syncFriendList()
	return result
}

func (f *Friend) setFriendRemark(params sdk_params_callback.SetFriendRemarkParams, callback common.Base, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.SetFriendRemarkReq{}
	apiReq.OperationID = operationID
	apiReq.ToUserID = params.ToUserID
	apiReq.FromUserID = f.loginUserID
	result := f.p.PostFatalCallback(callback, constant.SetFriendComment, apiReq, f.token)
	f.syncFriendList()
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

//
//func (u *UserRelated) doFriendList() {
//	friendsInfoOnServer, err := u.getServerFriendList()
//	if err != nil {
//		return
//	}
//	friendsInfoOnServerInterface := make([]diff, 0)
//	for _, v := range friendsInfoOnServer {
//		friendsInfoOnServerInterface = append(friendsInfoOnServerInterface, v)
//	}
//	friendsInfoOnLocal, err := u.getLocalFriendList()
//	if err != nil {
//		return
//	}
//	friendsInfoOnLocalInterface := make([]diff, 0)
//	for _, v := range friendsInfoOnLocal {
//		friendsInfoOnLocalInterface = append(friendsInfoOnLocalInterface, v)
//	}
//	aInBNot, bInANot, sameA, _ := checkDiff(friendsInfoOnServerInterface, friendsInfoOnLocalInterface)
//	if len(aInBNot) > 0 {
//		for _, index := range aInBNot {
//			if friendInfoStruct, ok := friendsInfoOnServerInterface[index].Value().(friendInfo); ok {
//				err = u.insertIntoTheFriendToFriendInfo(friendInfoStruct.UID, friendInfoStruct.Name, friendInfoStruct.Comment, friendInfoStruct.Icon, friendInfoStruct.Gender, friendInfoStruct.Mobile, friendInfoStruct.Birth, friendInfoStruct.Email, friendInfoStruct.Ex)
//				if err != nil {
//					sdkLog(err.Error())
//					return
//				}
//				jsonFriendInfo, _ := json.Marshal(friendInfoStruct)
//				u.friendListener.OnFriendListAdded(string(jsonFriendInfo))
//			}
//		}
//	}
//
//	if len(bInANot) > 0 {
//		for _, index := range bInANot {
//			err = u.delTheFriendFromFriendInfo(friendsInfoOnLocalInterface[index].Key())
//			if err != nil {
//				sdkLog(err.Error())
//				return
//			}
//			jsonFriendInfo, _ := json.Marshal(friendsInfoOnLocal[index])
//			u.friendListener.OnFriendListDeleted(string(jsonFriendInfo))
//			//_ = u.triggerCmdDeleteConversationAndMessage(friendsInfoOnLocalInterface[index].Key(), GetConversationIDBySessionType(friendsInfoOnLocalInterface[index].Key(), SingleChatType), SingleChatType)
//		}
//	}
//
//	if len(sameA) > 0 {
//		for _, index := range sameA {
//			if friendInfoStruct, ok := friendsInfoOnServerInterface[index].Value().(friendInfo); ok {
//				err = u.updateTheFriendInfo(friendInfoStruct.UID, friendInfoStruct.Name, friendInfoStruct.Comment, friendInfoStruct.Icon, friendInfoStruct.Gender, friendInfoStruct.Mobile, friendInfoStruct.Birth, friendInfoStruct.Email, friendInfoStruct.Ex)
//				if err != nil {
//					sdkLog(err.Error())
//					return
//				}
//				jsonFriendInfo, _ := json.Marshal(friendInfoStruct)
//				u.friendListener.OnFriendInfoChanged(string(jsonFriendInfo))
//			}
//		}
//	}
//}

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
	resp, err := network.Post2Api(constant.GetFriendListRouter, apiReq, f.token)
	commData, err := common.CheckErrAndRespReturn(err, resp, apiReq.OperationID)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}

	realData := server_api_params.GetFriendListResp{}
	err = mapstructure.Decode(commData.Data, &realData.FriendInfoList)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}

	return realData.FriendInfoList, nil
}

func (f *Friend) doBlackList() {
	//
	//blackListOnServer, err := u.getServerBlackList()
	//if err != nil {
	//	return
	//}
	//blackListOnServerInterface := make([]diff, 0)
	//for _, blackUser := range blackListOnServer {
	//	//blackListOnServerInterface = append(blackListOnServerInterface, blackUser)
	//}
	//
	//blackListOnLocal, err := u.getLocalBlackList()
	//if err != nil {
	//	return
	//}
	//blackListOnLocalInterface := make([]diff, 0)
	//for _, blackUser := range blackListOnLocal {z
	//	blackListOnLocalInterface = append(blackListOnLocalInterface, blackUser)
	//}
	//
	//aInBNot, bInANot, sameA, _ := checkDiff(blackListOnServerInterface, blackListOnLocalInterface)
	//
	//if len(aInBNot) > 0 {
	//	for _, index := range aInBNot {
	//		err = u.insertIntoTheUserToBlackList(blackListOnServer[index])
	//		if err != nil {
	//			sdkLog(err.Error())
	//			return
	//		}
	//		jsonAddBlackUserInfo, _ := json.Marshal(blackListOnServerInterface[index])
	//		u.friendListener.OnBlackListAdd(string(jsonAddBlackUserInfo))
	//	}
	//}
	//
	//if len(bInANot) > 0 {
	//	for _, index := range bInANot {
	//		err = u.delTheUserFromBlackList(blackListOnLocalInterface[index].Key())
	//		if err != nil {
	//			sdkLog(err.Error())
	//			return
	//		}
	//		jsonDelBlackUserInfo, _ := json.Marshal(blackListOnLocal[index])
	//		u.friendListener.OnBlackListDeleted(string(jsonDelBlackUserInfo))
	//	}
	//}
	//if len(bInANot) > 0 || len(aInBNot) > 0 {
	//	_ = triggerCmdFriend()
	//}
	//
	//if len(sameA) > 0 {
	//	for _, index := range sameA {
	//		//interface--->struct
	//		if blackListStruct, ok := blackListOnServerInterface[index].Value().(userInfo); ok {
	//			_ = u.updateBlackList(blackListStruct)
	//		}
	//	}
	//}
}

func (f *Friend) getServerBlackList(operationID string) ([]*server_api_params.PublicUserInfo, error) {
	apiReq := server_api_params.GetBlackListReq{OperationID: operationID, FromUserID: f.loginUserID}
	resp, err := network.Post2Api(constant.GetBlackListRouter, apiReq, f.token)
	commData, err := common.CheckErrAndRespReturn(err, resp, apiReq.OperationID)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	realData := server_api_params.GetBlackListResp{}
	mapstructure.Decode(commData.Data, &realData.BlackUserInfoList)
	return realData.BlackUserInfoList, nil
}

func (f *Friend) getServerFriendApplication(operationID string) ([]*server_api_params.FriendRequest, error) {
	apiReq := server_api_params.GetFriendApplyListReq{OperationID: operationID, FromUserID: f.loginUserID}
	resp, err := network.Post2Api(constant.GetFriendApplicationListRouter, apiReq, f.token)
	commData, err := common.CheckErrAndRespReturn(err, resp, apiReq.OperationID)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	realData := server_api_params.GetFriendApplyListResp{}
	mapstructure.Decode(commData.Data, &realData.FriendRequestList)
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
	result := f.p.PostFatalCallback(callback, constant.AddBlackListRouter, apiReq, f.token)
	f.syncBlackList()
	return result

}
func (f *Friend) removeBlack(callback common.Base, deleteUid, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.RemoveBlackListReq{}
	apiReq.ToUserID = deleteUid
	apiReq.FromUserID = f.loginUserID
	apiReq.OperationID = operationID
	result := f.p.PostFatalCallback(callback, constant.RemoveBlackListRouter, apiReq, f.token)
	f.syncBlackList()
	return result

}

//
//func (u *UserRelated) doAcceptOrRefuseApplicationCall(sendUid string, flag int32) {
//	sdkLog("doAcceptOrRefuseApplicationCall", sendUid, flag)
//
//	var ui2GetUserInfo ui2ClientCommonReq
//	ui2GetUserInfo.UidList = append(ui2GetUserInfo.UidList, sendUid)
//	resp, err := post2Api(getUserInfoRouter, paramsGetUserInfo{UidList: ui2GetUserInfo.UidList, OperationID: operationIDGenerator()}, u.token)
//	if err != nil {
//		sdkLog("getUserInfo failed", err)
//		return
//	}
//	var vgetUserInfoResp getUserInfoResp
//	err = json.Unmarshal(resp, &vgetUserInfoResp)
//	if err != nil {
//
//	}
//	if vgetUserInfoResp.ErrCode != 0 {
//		sdkLog(vgetUserInfoResp.ErrCode, vgetUserInfoResp.ErrMsg)
//		return
//	}
//	var appUserNode applyUserInfo
//	appUserNode.Uid = vgetUserInfoResp.Data[0].Uid
//	appUserNode.Name = vgetUserInfoResp.Data[0].Name
//	appUserNode.Icon = vgetUserInfoResp.Data[0].Icon
//	appUserNode.Gender = vgetUserInfoResp.Data[0].Gender
//	appUserNode.Mobile = vgetUserInfoResp.Data[0].Mobile
//	appUserNode.Birth = vgetUserInfoResp.Data[0].Birth
//	appUserNode.Email = vgetUserInfoResp.Data[0].Email
//	appUserNode.Ex = vgetUserInfoResp.Data[0].Ex
//	appUserNode.Flag = flag
//
//	jsonInfo, err := json.Marshal(appUserNode)
//	if err != nil {
//		sdkLog("doAcceptOrRefuseApplication json marshal failed")
//		return
//	}
//	sdkLog(flag)
//	if flag == 1 {
//		u.friendListener.OnFriendApplicationListAccept(string(jsonInfo))
//	}
//	if flag == -1 {
//		u.friendListener.OnFriendApplicationListReject(string(jsonInfo))
//	}
//}

func (f *Friend) syncSelfFriendApplication() {

}

func (f *Friend) syncFriendApplication() {

	svrList, err := f.getServerFriendApplication("")
	if err != nil {
		log.NewError("0", "getServerFriendList failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalFriendRequest(svrList)
	onLocal, err := f.db.GetRecvFriendApplication()
	if err != nil {
		log.NewError("0", "GetRecvFriendApplication failed ", err.Error())
		return
	}
	log.NewInfo("0", "svrList", svrList)
	log.NewInfo("0", "onServer", onServer)
	log.NewInfo("0", "onLocal", onLocal)

	aInBNot, bInANot, sameA, _ := common.CheckFriendRequestDiff(onServer, onLocal)
	for _, index := range aInBNot {
		err := f.db.InsertFriendRequest(onServer[index])
		if err != nil {
			log.NewError("0", "InsertFriendRequest failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := f.db.UpdateFriendRequest(onServer[index])
		if err != nil {
			log.NewError("0", "UpdateFriendRequest failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriendRequestBothUserID(onServer[index].FromUserID, onServer[index].ToUserID)
		if err != nil {
			log.NewError("0", "_deleteFriendRequestBothUserID failed ", err.Error())
			continue
		}
	}
}

func (f *Friend) syncFriendList() {
	svrList, err := f.getServerFriendList("")
	if err != nil {
		log.NewError("0", "getServerFriendList failed ", err.Error())
		return
	}
	log.NewInfo("0", "svrList", svrList)
	friendsInfoOnServer := common.TransferToLocalFriend(svrList)
	friendsInfoOnLocal, err := f.db.GetAllFriendList()
	if err != nil {
		log.NewError("0", "_getAllFriendList failed ", err.Error())
		return
	}

	log.NewInfo("0", "friendsInfoOnServer", friendsInfoOnServer)
	log.NewInfo("0", "friendsInfoOnLocal", friendsInfoOnLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckFriendListDiff(friendsInfoOnServer, friendsInfoOnLocal)
	log.NewInfo("0", "checkFriendListDiff", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := f.db.InsertFriend(friendsInfoOnServer[index])
		if err != nil {
			log.NewError("0", "_insertFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := f.db.UpdateFriend(friendsInfoOnServer[index])
		if err != nil {
			log.NewError("0", "_updateFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriend(friendsInfoOnLocal[index].FriendUserID)
		if err != nil {
			log.NewError("0", "_deleteFriend failed ", err.Error())
			continue
		}
	}
}

func (f *Friend) syncBlackList() {
	svrList, err := f.getServerBlackList("")
	if err != nil {
		log.NewError("0", "getServerBlackList failed ", err.Error())
		return
	}
	log.NewInfo("0", "svrList", svrList)
	blackListOnServer := common.TransferToLocalBlack(svrList, f.loginUserID)
	blackListOnLocal, err := f.db.GetBlackList()
	if err != nil {
		log.NewError("0", "_getBlackList failed ", err.Error())
		return
	}

	log.NewInfo("0", "blackListOnServer", blackListOnServer)
	log.NewInfo("0", "blackListOnlocal", blackListOnLocal)
	aInBNot, bInANot, sameA, _ := common.CheckBlackListDiff(blackListOnServer, blackListOnLocal)
	for _, index := range aInBNot {
		err := f.db.InsertBlack(blackListOnServer[index])
		if err != nil {
			log.NewError("0", "_insertFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := f.db.UpdateBlack(blackListOnServer[index])
		if err != nil {
			log.NewError("0", "_updateFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteBlack(blackListOnLocal[index].BlockUserID)
		if err != nil {
			log.NewError("0", "_deleteFriend failed ", err.Error())
			continue
		}
	}

}

func (u *Friend) addFriendNew(msg *server_api_params.MsgData) {
	//utils.sdkLog("addFriend start ")
	//u.syncFriendApplication()
	//
	//var ui2GetUserInfo open_im_sdk.ui2ClientCommonReq
	//ui2GetUserInfo.UidList = append(ui2GetUserInfo.UidList, msg.SendID)
	//resp, err := utils.post2Api(open_im_sdk.getUserInfoRouter, open_im_sdk.paramsGetUserInfo{UidList: ui2GetUserInfo.UidList, OperationID: utils.operationIDGenerator()}, u.token)
	//if err != nil {
	//	utils.sdkLog("getUserInfo failed", err)
	//	return
	//}
	//var vgetUserInfoResp open_im_sdk.getUserInfoResp
	//err = json.Unmarshal(resp, &vgetUserInfoResp)
	//if err != nil {
	//	utils.sdkLog("Unmarshal failed, ", err.Error())
	//	return
	//}
	//if vgetUserInfoResp.ErrCode != 0 {
	//	utils.sdkLog(vgetUserInfoResp.ErrCode, vgetUserInfoResp.ErrMsg)
	//	return
	//}
	//if len(vgetUserInfoResp.Data) == 0 {
	//	utils.sdkLog(vgetUserInfoResp.ErrCode, vgetUserInfoResp.ErrMsg, msg)
	//	return
	//}
	//var appUserNode open_im_sdk.applyUserInfo
	//appUserNode.Uid = vgetUserInfoResp.Data[0].Uid
	//appUserNode.Name = vgetUserInfoResp.Data[0].Name
	//appUserNode.Icon = vgetUserInfoResp.Data[0].Icon
	//appUserNode.Gender = vgetUserInfoResp.Data[0].Gender
	//appUserNode.Mobile = vgetUserInfoResp.Data[0].Mobile
	//appUserNode.Birth = vgetUserInfoResp.Data[0].Birth
	//appUserNode.Email = vgetUserInfoResp.Data[0].Email
	//appUserNode.Ex = vgetUserInfoResp.Data[0].Ex
	//appUserNode.Flag = 0
	//
	//jsonInfo, err := json.Marshal(appUserNode)
	//if err != nil {
	//	utils.sdkLog("  marshal failed", err.Error())
	//	return
	//}
	//u.friendListener.OnFriendApplicationListAdded(string(jsonInfo))
}

func (u *Friend) DoFriendMsg(msg *server_api_params.MsgData) {
	//utils.sdkLog("doFriendMsg ", msg)
	//if u.cb == nil || u.friendListener == nil {
	//	utils.sdkLog("listener is null")
	//	return
	//}
	//
	//if msg.SendID == u.loginUserID && msg.SenderPlatformID == u.SvrConf.Platform {
	//	utils.sdkLog("sync msg ", msg.ContentType)
	//	return
	//}
	//
	//go func() {
	//	switch msg.ContentType {
	//	case constant.AddFriendTip:
	//		utils.sdkLog("addFriendNew ", msg)
	//		u.addFriendNew(msg) //
	//	case constant.AcceptFriendApplicationTip:
	//		utils.sdkLog("acceptFriendApplicationNew ", msg)
	//		u.acceptFriendApplicationNew(msg)
	//	case constant.RefuseFriendApplicationTip:
	//		utils.sdkLog("refuseFriendApplicationNew ", msg)
	//		u.refuseFriendApplicationNew(msg)
	//	case constant.SetSelfInfoTip:
	//		utils.sdkLog("setSelfInfo ", msg)
	//		u.setSelfInfo(msg)
	//		//	case KickOnlineTip:
	//		//		sdkLog("kickOnline ", msg)
	//		//		u.kickOnline(&msg)
	//	default:
	//		utils.sdkLog("type failed, ", msg)
	//	}
	//}()
}

func (u *Friend) acceptFriendApplicationNew(msg *server_api_params.MsgData) {
	//utils.LogBegin(msg.ContentType, msg.ServerMsgID, msg.ClientMsgID)
	//u.syncFriendList()
	//utils.sdkLog(msg.SendID, msg.RecvID)
	//utils.sdkLog("acceptFriendApplicationNew", msg.ServerMsgID, msg)
	//
	//fInfoList, err := u.getServerFriendList()
	//if err != nil {
	//	return
	//}
	//utils.sdkLog("fInfoList", fInfoList)

	//for _, fInfo := range fInfoList {
	//	if fInfo.UID == msg.SendID {
	//		jData, err := json.Marshal(fInfo)
	//		if err != nil {
	//			sdkLog("err: ", err.Error())
	//			return
	//		}
	//		u.friendListener.OnFriendListAdded(string(jData))
	//		u.friendListener.OnFriendApplicationListAccept(string(jData))
	//		return
	//	}
	//}
}

func (u *Friend) refuseFriendApplicationNew(msg *server_api_params.MsgData) {
	//applyList, err := u.getServerSelfApplication()
	//
	//if err != nil {
	//	return
	//}
	//for _, v := range applyList {
	//	if v.Uid == msg.SendID {
	//		jData, err := json.Marshal(v)
	//		if err != nil {
	//			return
	//		}
	//		u.friendListener.OnFriendApplicationListReject(string(jData))
	//		return
	//	}
	//}

}
