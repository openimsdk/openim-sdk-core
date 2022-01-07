package friend

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/log"
	"open_im_sdk/open_im_sdk/sdk_params_callback"
	"open_im_sdk/open_im_sdk/server_api_params"
	"open_im_sdk/open_im_sdk/utils"
)

type FriendListener struct {
	friendListener open_im_sdk.OnFriendshipListener
}

func (u *open_im_sdk.UserRelated) getDesignatedFriendsInfo(callback open_im_sdk.Base, friendUserIDList sdk_params_callback.GetDesignatedFriendsInfoParams, operationID string) sdk_params_callback.GetDesignatedFriendsInfoCallback {
	log.NewInfo(operationID, utils.GetSelfFuncName(), friendUserIDList)
	blackList, err := u._getBlackInfoList(friendUserIDList)
	utils.checkErr(callback, err, operationID)
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
	localFriendList, err := u._getFriendInfoList(pureFriendUserIDList)
	utils.checkErr(callback, err, operationID)
	log.NewInfo(operationID, "_getFriendInfoList ", pureFriendUserIDList, localFriendList)
	return localFriendList
}

func (u *open_im_sdk.UserRelated) addFriend(callback open_im_sdk.Base, addFriendParams sdk_params_callback.AddFriendParams, operationID string) *server_api_params.CommDataResp {
	log.NewInfo(operationID, "addFriend args: ", addFriendParams)
	apiReq := server_api_params.AddFriendReq{}
	apiReq.ToUserID = addFriendParams.ToUserID
	apiReq.FromUserID = u.loginUserID
	apiReq.ReqMsg = addFriendParams.ReqMsg
	apiReq.OperationID = operationID
	resp, err := utils.post2Api(open_im_sdk.addFriendRouter, apiReq, u.token)
	log.NewInfo(apiReq.OperationID, "post2Api ", open_im_sdk.addFriendRouter, apiReq, string(resp))
	return utils.checkErrAndResp(callback, err, resp, operationID)
}

func (u *open_im_sdk.UserRelated) getRecvFriendApplicationList(callback open_im_sdk.Base, operationID string) sdk_params_callback.GetRecvFriendApplicationListCallback {
	log.NewInfo(operationID, "getRecvFriendApplicationList args: ")
	friendApplicationList, err := u._getRecvFriendApplication()
	utils.checkErr(callback, err, operationID)
	return friendApplicationList
}

func (u *open_im_sdk.UserRelated) getSendFriendApplicationList(callback open_im_sdk.Base, operationID string) sdk_params_callback.GetSendFriendApplicationListCallback {
	return nil
}

func (u *open_im_sdk.UserRelated) processFriendApplication(callback open_im_sdk.Base, params sdk_params_callback.ProcessFriendApplicationParams, handleResult int32, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.AddFriendResponseReq{}
	apiReq.FromUserID = u.loginUserID
	apiReq.ToUserID = params.ToUserID
	apiReq.Flag = handleResult
	apiReq.OperationID = operationID
	apiReq.HandleMsg = params.HandleMsg
	resp, err := utils.post2Api(open_im_sdk.addFriendResponse, apiReq, u.token)
	r := utils.checkErrAndResp(callback, err, resp, operationID)
	u.syncFriendApplication()
	return r
}

func (u *open_im_sdk.UserRelated) checkFriend(callback open_im_sdk.Base, userIDList sdk_params_callback.CheckFriendParams, operationID string) sdk_params_callback.CheckFriendCallback {
	friendList, err := u._getFriendInfoList(userIDList)
	utils.checkErr(callback, err, operationID)
	blackList, err := u._getBlackInfoList(userIDList)
	utils.checkErr(callback, err, operationID)
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

func (u *open_im_sdk.UserRelated) deleteFriend(FriendUserID string, callback open_im_sdk.Base, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.DeleteFriendReq{}
	apiReq.ToUserID = FriendUserID
	apiReq.FromUserID = u.loginUserID
	apiReq.OperationID = operationID
	resp, err := utils.post2Api(open_im_sdk.deleteFriendRouter, apiReq, u.token)
	result := utils.checkErrAndResp(callback, err, resp, operationID)
	u.syncFriendList()
	return result
}

func (u *open_im_sdk.UserRelated) setFriendRemark(params sdk_params_callback.SetFriendRemarkParams, callback open_im_sdk.Base, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.SetFriendRemarkReq{}
	apiReq.OperationID = operationID
	apiReq.ToUserID = params.ToUserID
	apiReq.FromUserID = u.loginUserID
	resp, err := utils.post2Api(open_im_sdk.setFriendComment, apiReq, u.token)
	result := utils.checkErrAndResp(callback, err, resp, operationID)
	u.syncFriendList()
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

func (u *open_im_sdk.UserRelated) getLocalFriendList() ([]open_im_sdk.friendInfo, error) {
	//Take out the friend list and judge whether it is in the blacklist again to prevent nested locks
	localFriendList, err := u.getLocalFriendList22()
	if err != nil {
		return nil, err
	}
	for index, v := range localFriendList {
		//check friend is in blacklist
		blackUser, err := u.getBlackUsInfoByUid(v.UID)
		if err != nil {
			utils.sdkLog(err.Error())
		}
		if blackUser.Uid != "" {
			localFriendList[index].IsInBlackList = 1
		}
	}
	return localFriendList, nil
}

func (u *open_im_sdk.UserRelated) getServerFriendList() ([]*open_im_sdk.FriendInfo, error) {
	apiReq := server_api_params.GetFriendListReq{OperationID: utils.operationIDGenerator(), FromUserID: u.loginUserID}
	resp, err := utils.post2Api(open_im_sdk.getFriendListRouter, apiReq, u.token)
	commData, err := utils.checkErrAndRespReturn(err, resp, apiReq.OperationID)
	if err != nil {
		return nil, utils.wrap(err, apiReq.OperationID)
	}

	realData := server_api_params.GetFriendListResp{}
	err = mapstructure.Decode(commData.Data, &realData.FriendInfoList)
	if err != nil {
		return nil, utils.wrap(err, apiReq.OperationID)
	}

	return realData.FriendInfoList, nil
}

func (u *open_im_sdk.UserRelated) doBlackList() {
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

func (u *open_im_sdk.UserRelated) getServerBlackList() ([]*open_im_sdk.PublicUserInfo, error) {
	apiReq := server_api_params.GetBlackListReq{OperationID: utils.operationIDGenerator(), FromUserID: u.loginUserID}
	resp, err := utils.post2Api(open_im_sdk.getBlackListRouter, apiReq, u.token)
	commData, err := utils.checkErrAndRespReturn(err, resp, apiReq.OperationID)
	if err != nil {
		return nil, utils.wrap(err, apiReq.OperationID)
	}
	realData := server_api_params.GetBlackListResp{}
	mapstructure.Decode(commData.Data, &realData.BlackUserInfoList)
	return realData.BlackUserInfoList, nil
}

func (u *open_im_sdk.UserRelated) getServerFriendApplication() ([]*open_im_sdk.FriendRequest, error) {
	apiReq := server_api_params.GetFriendApplyListReq{OperationID: utils.operationIDGenerator(), FromUserID: u.loginUserID}
	resp, err := utils.post2Api(open_im_sdk.getFriendApplicationListRouter, apiReq, u.token)
	commData, err := utils.checkErrAndRespReturn(err, resp, apiReq.OperationID)
	if err != nil {
		return nil, utils.wrap(err, apiReq.OperationID)
	}

	realData := server_api_params.GetFriendApplyListResp{}
	mapstructure.Decode(commData.Data, &realData.FriendRequestList)
	return realData.FriendRequestList, nil
}

func (u *open_im_sdk.UserRelated) getServerSelfApplication() ([]open_im_sdk.applyUserInfo, error) {
	resp, err := utils.post2Api(open_im_sdk.getSelfApplicationListRouter, open_im_sdk.paramsCommonReq{OperationID: utils.operationIDGenerator()}, u.token)
	if err != nil {
		return nil, err
	}
	var vgetFriendApplyListResp open_im_sdk.getFriendApplyListResp
	err = json.Unmarshal(resp, &vgetFriendApplyListResp)
	if err != nil {
		utils.sdkLog("unmarshal failed, ", err.Error())
		return nil, err
	}
	if vgetFriendApplyListResp.ErrCode != 0 {
		utils.sdkLog("errcode: ", vgetFriendApplyListResp.ErrCode, "errmsg: ", vgetFriendApplyListResp.ErrMsg)
		return nil, err
	}
	return vgetFriendApplyListResp.Data, nil
}
func (u *open_im_sdk.UserRelated) addBlack(callback open_im_sdk.Base, blackUid, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.AddBlacklistReq{}
	apiReq.ToUserID = blackUid
	apiReq.FromUserID = u.loginUserID
	apiReq.OperationID = operationID
	resp, err := utils.post2Api(open_im_sdk.addBlackListRouter, apiReq, u.token)
	r := utils.checkErrAndResp(callback, err, resp, operationID)
	u.syncBlackList()
	return r

}
func (u *open_im_sdk.UserRelated) removeBlack(callback open_im_sdk.Base, deleteUid, operationID string) *server_api_params.CommDataResp {
	apiReq := server_api_params.RemoveBlackListReq{}
	apiReq.ToUserID = deleteUid
	apiReq.FromUserID = u.loginUserID
	apiReq.OperationID = operationID
	resp, err := utils.post2Api(open_im_sdk.removeBlackListRouter, apiReq, u.token)
	r := utils.checkErrAndResp(callback, err, resp, operationID)
	u.syncBlackList()
	return r

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

func (u *open_im_sdk.UserRelated) syncSelfFriendApplication() {

}

func (u *open_im_sdk.UserRelated) syncFriendApplication() {
	svrList, err := u.getServerFriendApplication()
	if err != nil {
		log.NewError("0", "getServerFriendList failed ", err.Error())
		return
	}
	onServer := utils.transferToLocalFriendRequest(svrList)
	onLocal, err := u._getRecvFriendApplication()
	if err != nil {
		log.NewError("0", "_getAllFriendList failed ", err.Error())
		return
	}
	log.NewInfo("0", "svrList", svrList)
	log.NewInfo("0", "onServer", onServer)
	log.NewInfo("0", "onLocal", onLocal)

	aInBNot, bInANot, sameA, _ := utils.checkFriendRequestDiff(onServer, onLocal)
	for _, index := range aInBNot {
		err := u._insertFriendRequest(onServer[index])
		if err != nil {
			log.NewError("0", "_insertFriendRequest failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := u._updateFriendRequest(onServer[index])
		if err != nil {
			log.NewError("0", "_updateFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := u._deleteFriendRequestBothUserID(onServer[index].FromUserID, onServer[index].ToUserID)
		if err != nil {
			log.NewError("0", "_deleteFriendRequestBothUserID failed ", err.Error())
			continue
		}
	}
}

func (u *open_im_sdk.UserRelated) syncFriendList() {
	svrList, err := u.getServerFriendList()
	if err != nil {
		log.NewError("0", "getServerFriendList failed ", err.Error())
		return
	}
	log.NewInfo("0", "svrList", svrList)
	friendsInfoOnServer := utils.transferToLocalFriend(svrList)
	friendsInfoOnLocal, err := u._getAllFriendList()
	if err != nil {
		log.NewError("0", "_getAllFriendList failed ", err.Error())
		return
	}

	log.NewInfo("0", "friendsInfoOnServer", friendsInfoOnServer)
	log.NewInfo("0", "friendsInfoOnLocal", friendsInfoOnLocal)
	aInBNot, bInANot, sameA, sameB := utils.checkFriendListDiff(friendsInfoOnServer, friendsInfoOnLocal)
	log.NewInfo("0", "checkFriendListDiff", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := u._insertFriend(friendsInfoOnServer[index])
		if err != nil {
			log.NewError("0", "_insertFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := u._updateFriend(friendsInfoOnServer[index])
		if err != nil {
			log.NewError("0", "_updateFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := u._deleteFriend(friendsInfoOnLocal[index].FriendUserID)
		if err != nil {
			log.NewError("0", "_deleteFriend failed ", err.Error())
			continue
		}
	}
}

func (u *open_im_sdk.UserRelated) syncBlackList() {
	svrList, err := u.getServerBlackList()
	if err != nil {
		log.NewError("0", "getServerBlackList failed ", err.Error())
		return
	}
	log.NewInfo("0", "svrList", svrList)
	blackListOnServer := utils.transferToLocalBlack(svrList, u.loginUserID)
	blackListOnLocal, err := u._getBlackList()
	if err != nil {
		log.NewError("0", "_getBlackList failed ", err.Error())
		return
	}

	log.NewInfo("0", "blackListOnServer", blackListOnServer)
	log.NewInfo("0", "blackListOnlocal", blackListOnLocal)
	aInBNot, bInANot, sameA, _ := utils.checkBlackListDiff(blackListOnServer, blackListOnLocal)
	for _, index := range aInBNot {
		err := u._insertBlack(blackListOnServer[index])
		if err != nil {
			log.NewError("0", "_insertFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := u._updateBlack(blackListOnServer[index])
		if err != nil {
			log.NewError("0", "_updateFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := u._deleteBlack(blackListOnLocal[index].BlockUserID)
		if err != nil {
			log.NewError("0", "_deleteFriend failed ", err.Error())
			continue
		}
	}

}
