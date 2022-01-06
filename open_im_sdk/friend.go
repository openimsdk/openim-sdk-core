package open_im_sdk

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
)

type FriendListener struct {
	friendListener OnFriendshipListener
}

func (u *UserRelated) getDesignatedFriendsInfo(callback Base, friendUserIDList GetDesignatedFriendsInfoParams, operationID string) GetDesignatedFriendsInfoCallback {
	blackList, err := u._getBlackInfoList(friendUserIDList)
	checkErr(callback, err, operationID)

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
	checkErr(callback, err, operationID)

	return localFriendList
}

func (u *UserRelated) addFriend(callback Base, addFriendParams AddFriendParams, operationID string) *CommDataResp {
	apiReq := AddFriendReq{}
	apiReq.ToUserID = addFriendParams.ToUserID
	apiReq.FromUserID = u.loginUserID
	apiReq.ReqMsg = addFriendParams.ReqMsg
	apiReq.OperationID = operationID
	resp, err := post2Api(addFriendRouter, apiReq, u.token)
	NewInfo(apiReq.OperationID, "post2Api ", addFriendRouter, apiReq)
	return checkErrAndResp(callback, err, resp, operationID)
}

func (u *UserRelated) getRecvFriendApplicationList(callback Base, operationID string) GetRecvFriendApplicationListCallback {
	friendApplicationList, err := u._getRecvFriendApplication()
	checkErr(callback, err, operationID)
	return friendApplicationList
}

func (u *UserRelated) getSendFriendApplicationList(callback Base, operationID string) GetSendFriendApplicationListCallback {
	return nil
}

func (u *UserRelated) processFriendApplication(callback Base, params ProcessFriendApplicationParams, handleResult int32, operationID string) *CommDataResp {
	apiReq := AddFriendResponseReq{}
	apiReq.FromUserID = u.loginUserID
	apiReq.ToUserID = params.ToUserID
	apiReq.Flag = handleResult
	apiReq.OperationID = operationID
	apiReq.HandleMsg = params.HandleMsg
	resp, err := post2Api(addFriendResponse, apiReq, u.token)
	r := checkErrAndResp(callback, err, resp, operationID)
	u.syncFriendApplication()
	return r
}

func (u *UserRelated) checkFriend(callback Base, userIDList CheckFriendParams, operationID string) CheckFriendCallback {
	friendList, err := u._getFriendInfoList(userIDList)
	checkErr(callback, err, operationID)
	blackList, err := u._getBlackInfoList(userIDList)
	checkErr(callback, err, operationID)
	var checkFriendCallback CheckFriendCallback
	for _, v := range userIDList {
		var r UserIDResult
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

func (u *UserRelated) deleteFriend(FriendUserID string, callback Base, operationID string) *CommDataResp {
	apiReq := DeleteFriendReq{}
	apiReq.ToUserID = FriendUserID
	apiReq.FromUserID = u.loginUserID
	resp, err := post2Api(deleteFriendRouter, apiReq, u.token)
	result := checkErrAndResp(callback, err, resp, operationID)
	u.syncFriendList()
	return result
}

func (u *UserRelated) setFriendRemark(params SetFriendRemarkParams, callback Base, operationID string) *CommDataResp {
	apiReq := SetFriendRemarkReq{}
	apiReq.OperationID = operationID
	apiReq.ToUserID = params.ToUserID
	apiReq.FromUserID = u.loginUserID
	resp, err := post2Api(setFriendComment, apiReq, u.token)
	result := checkErrAndResp(callback, err, resp, operationID)
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

func (u *UserRelated) getLocalFriendList() ([]friendInfo, error) {
	//Take out the friend list and judge whether it is in the blacklist again to prevent nested locks
	localFriendList, err := u.getLocalFriendList22()
	if err != nil {
		return nil, err
	}
	for index, v := range localFriendList {
		//check friend is in blacklist
		blackUser, err := u.getBlackUsInfoByUid(v.UID)
		if err != nil {
			sdkLog(err.Error())
		}
		if blackUser.Uid != "" {
			localFriendList[index].IsInBlackList = 1
		}
	}
	return localFriendList, nil
}

func (u *UserRelated) getServerFriendList() ([]*FriendInfo, error) {
	apiReq := GetFriendListReq{OperationID: operationIDGenerator(), FromUserID: u.loginUserID}
	resp, err := post2Api(getFriendListRouter, apiReq, u.token)
	commData, err := checkErrAndRespReturn(err, resp, apiReq.OperationID)
	if err != nil {
		return nil, wrap(err, apiReq.OperationID)
	}

	realData := GetFriendListResp{}
	mapstructure.Decode(commData.Data, &realData.FriendInfoList)
	return realData.FriendInfoList, nil
}

func (u *UserRelated) doBlackList() {

	blackListOnServer, err := u.getServerBlackList()
	if err != nil {
		return
	}
	blackListOnServerInterface := make([]diff, 0)
	for _, blackUser := range blackListOnServer {
		blackListOnServerInterface = append(blackListOnServerInterface, blackUser)
	}

	blackListOnLocal, err := u.getLocalBlackList()
	if err != nil {
		return
	}
	blackListOnLocalInterface := make([]diff, 0)
	for _, blackUser := range blackListOnLocal {
		blackListOnLocalInterface = append(blackListOnLocalInterface, blackUser)
	}

	aInBNot, bInANot, sameA, _ := checkDiff(blackListOnServerInterface, blackListOnLocalInterface)

	if len(aInBNot) > 0 {
		for _, index := range aInBNot {
			err = u.insertIntoTheUserToBlackList(blackListOnServer[index])
			if err != nil {
				sdkLog(err.Error())
				return
			}
			jsonAddBlackUserInfo, _ := json.Marshal(blackListOnServerInterface[index])
			u.friendListener.OnBlackListAdd(string(jsonAddBlackUserInfo))
		}
	}

	if len(bInANot) > 0 {
		for _, index := range bInANot {
			err = u.delTheUserFromBlackList(blackListOnLocalInterface[index].Key())
			if err != nil {
				sdkLog(err.Error())
				return
			}
			jsonDelBlackUserInfo, _ := json.Marshal(blackListOnLocal[index])
			u.friendListener.OnBlackListDeleted(string(jsonDelBlackUserInfo))
		}
	}
	if len(bInANot) > 0 || len(aInBNot) > 0 {
		_ = triggerCmdFriend()
	}

	if len(sameA) > 0 {
		for _, index := range sameA {
			//interface--->struct
			if blackListStruct, ok := blackListOnServerInterface[index].Value().(userInfo); ok {
				_ = u.updateBlackList(blackListStruct)
			}
		}
	}
}

func (u *UserRelated) getServerBlackList() ([]userInfo, error) {
	resp, err := post2Api(getBlackListRouter, paramsCommonReq{OperationID: operationIDGenerator()}, u.token)
	if err != nil {
		return nil, err
	}
	var vgetBlackListResp getBlackListResp
	err = json.Unmarshal(resp, &vgetBlackListResp)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return nil, err
	}
	if vgetBlackListResp.ErrCode != 0 {
		sdkLog("errcode: ", vgetBlackListResp.ErrCode, "errmsg: ", vgetBlackListResp.ErrMsg)
		return nil, err
	}
	return vgetBlackListResp.Data, nil
}

func (u *UserRelated) getServerFriendApplication() ([]*FriendRequest, error) {
	apiReq := GetFriendApplyListReq{OperationID: operationIDGenerator(), FromUserID: u.loginUserID}
	resp, err := post2Api(getFriendApplicationListRouter, apiReq, u.token)
	commData, err := checkErrAndRespReturn(err, resp, apiReq.OperationID)
	if err != nil {
		return nil, wrap(err, apiReq.OperationID)
	}

	realData := GetFriendApplyListResp{}
	mapstructure.Decode(commData.Data, &realData.FriendRequestList)
	return realData.FriendRequestList, nil
}

func (u *UserRelated) getServerSelfApplication() ([]applyUserInfo, error) {
	resp, err := post2Api(getSelfApplicationListRouter, paramsCommonReq{OperationID: operationIDGenerator()}, u.token)
	if err != nil {
		return nil, err
	}
	var vgetFriendApplyListResp getFriendApplyListResp
	err = json.Unmarshal(resp, &vgetFriendApplyListResp)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return nil, err
	}
	if vgetFriendApplyListResp.ErrCode != 0 {
		sdkLog("errcode: ", vgetFriendApplyListResp.ErrCode, "errmsg: ", vgetFriendApplyListResp.ErrMsg)
		return nil, err
	}
	return vgetFriendApplyListResp.Data, nil
}
func (u *UserRelated) addBlack(callback Base, blackUid, operationID string) *CommDataResp {
	apiReq := AddBlacklistReq{}
	apiReq.ToUserID = blackUid
	apiReq.FromUserID = u.loginUserID
	apiReq.OperationID = operationID
	resp, err := post2Api(addBlackListRouter, apiReq, u.token)
	r := checkErrAndResp(callback, err, resp, operationID)
	u.syncBlackList()
	return r

}
func (u *UserRelated) removeBlack(callback Base, deleteUid, operationID string) *CommDataResp {
	apiReq := RemoveBlackListReq{}
	apiReq.ToUserID = deleteUid
	apiReq.FromUserID = u.loginUserID
	apiReq.OperationID = operationID
	resp, err := post2Api(removeBlackListRouter, apiReq, u.token)
	r := checkErrAndResp(callback, err, resp, operationID)
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

func (u *UserRelated) syncSelfFriendApplication() {

}

func (u *UserRelated) syncFriendApplication() {
	svrList, err := u.getServerFriendApplication()
	if err != nil {
		NewError("0", "getServerFriendList failed ", err.Error())
		return
	}
	onServer := transferToLocalFriendRequest(svrList)
	onLocal, err := u._getSendFriendApplication()
	if err != nil {
		NewError("0", "_getAllFriendList failed ", err.Error())
		return
	}
	NewInfo("0", "svrList", svrList)
	NewInfo("0", "onServer", onServer)
	NewInfo("0", "onLocal", onLocal)

	aInBNot, bInANot, sameA, _ := checkFriendRequestDiff(onServer, onLocal)
	for _, index := range aInBNot {
		err := u._insertFriendRequest(onServer[index])
		if err != nil {
			NewError("0", "_insertFriendRequest failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := u._updateFriendRequest(onServer[index])
		if err != nil {
			NewError("0", "_updateFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := u._deleteFriendRequestBothUserID(onServer[index].FromUserID, onServer[index].ToUserID)
		if err != nil {
			NewError("0", "_deleteFriendRequestBothUserID failed ", err.Error())
			continue
		}
	}
}

func (u *UserRelated) syncFriendList() {
	svrList, err := u.getServerFriendList()
	if err != nil {
		NewError("0", "getServerFriendList failed ", err.Error())
		return
	}
	friendsInfoOnServer := transferToLocalFriend(svrList)
	friendsInfoOnLocal, err := u._getAllFriendList()
	if err != nil {
		NewError("0", "_getAllFriendList failed ", err.Error())
		return
	}
	NewInfo("0", "svrList", svrList)
	NewInfo("0", "friendsInfoOnServer", friendsInfoOnServer)
	NewInfo("0", "friendsInfoOnLocal", friendsInfoOnLocal)
	aInBNot, bInANot, sameA, _ := checkFriendListDiff(friendsInfoOnServer, friendsInfoOnLocal)
	for _, index := range aInBNot {
		err := u._insertFriend(friendsInfoOnServer[index])
		if err != nil {
			NewError("0", "_insertFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := u._updateFriend(friendsInfoOnServer[index])
		if err != nil {
			NewError("0", "_updateFriend failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := u._deleteFriend(friendsInfoOnLocal[index].FriendUserID)
		if err != nil {
			NewError("0", "_deleteFriend failed ", err.Error())
			continue
		}
	}
}

func (u *UserRelated) syncBlackList() {

	blackListOnServer, err := u.getServerBlackList()
	if err != nil {
		return
	}
	blackListOnServerInterface := make([]diff, 0)
	for _, blackUser := range blackListOnServer {
		blackListOnServerInterface = append(blackListOnServerInterface, blackUser)
	}

	blackListOnLocal, err := u.getLocalBlackList()
	if err != nil {
		return
	}
	blackListOnLocalInterface := make([]diff, 0)
	for _, blackUser := range blackListOnLocal {
		blackListOnLocalInterface = append(blackListOnLocalInterface, blackUser)
	}

	aInBNot, bInANot, sameA, _ := checkDiff(blackListOnServerInterface, blackListOnLocalInterface)

	if len(aInBNot) > 0 {
		for _, index := range aInBNot {
			err = u.insertIntoTheUserToBlackList(blackListOnServer[index])
			if err != nil {
				sdkLog(err.Error())
				return
			}

		}
	}

	if len(bInANot) > 0 {
		for _, index := range bInANot {
			err = u.delTheUserFromBlackList(blackListOnLocalInterface[index].Key())
			if err != nil {
				sdkLog(err.Error())
				return
			}

		}
	}

	if len(sameA) > 0 {
		for _, index := range sameA {
			//interface--->struct
			if blackListStruct, ok := blackListOnServerInterface[index].Value().(userInfo); ok {
				_ = u.updateBlackList(blackListStruct)
			}
		}
	}
}
