package open_im_sdk

import (
	"encoding/json"
	"errors"
	"time"
)

type Friend struct {
	friendListener OnFriendshipListener
	//ch             chan cmd2Value
}

func sendCmd(ch chan cmd2Value, value cmd2Value, timeout int64) error {
	var flag = 0
	select {
	case ch <- value:
		flag = 1
	case <-time.After(time.Second * time.Duration(timeout)):
		flag = 2
	}
	if flag == 1 {
		return nil
	} else {
		sdkLog("send cmd timeout, ", timeout, value)
		return errors.New("send cmd timeout")
	}
}

func (u *UserRelated) doFriendList() {
	friendsInfoOnServer, err := u.getServerFriendList()
	if err != nil {
		return
	}
	friendsInfoOnServerInterface := make([]diff, 0)
	for _, v := range friendsInfoOnServer {
		friendsInfoOnServerInterface = append(friendsInfoOnServerInterface, v)
	}
	friendsInfoOnLocal, err := u.getLocalFriendList()
	if err != nil {
		return
	}
	friendsInfoOnLocalInterface := make([]diff, 0)
	for _, v := range friendsInfoOnLocal {
		friendsInfoOnLocalInterface = append(friendsInfoOnLocalInterface, v)
	}
	aInBNot, bInANot, sameA, _ := checkDiff(friendsInfoOnServerInterface, friendsInfoOnLocalInterface)
	if len(aInBNot) > 0 {
		for _, index := range aInBNot {
			if friendInfoStruct, ok := friendsInfoOnServerInterface[index].Value().(friendInfo); ok {
				err = u.insertIntoTheFriendToFriendInfo(friendInfoStruct.UID, friendInfoStruct.Name, friendInfoStruct.Comment, friendInfoStruct.Icon, friendInfoStruct.Gender, friendInfoStruct.Mobile, friendInfoStruct.Birth, friendInfoStruct.Email, friendInfoStruct.Ex)
				if err != nil {
					log(err.Error())
					return
				}
				jsonFriendInfo, _ := json.Marshal(friendInfoStruct)
				u.friendListener.OnFriendListAdded(string(jsonFriendInfo))
			}
		}
	}

	if len(bInANot) > 0 {
		for _, index := range bInANot {
			err = u.delTheFriendFromFriendInfo(friendsInfoOnLocalInterface[index].Key())
			if err != nil {
				log(err.Error())
				return
			}
			jsonFriendInfo, _ := json.Marshal(friendsInfoOnLocal[index])
			u.friendListener.OnFriendListDeleted(string(jsonFriendInfo))
			_ = u.triggerCmdDeleteConversationAndMessage(friendsInfoOnLocalInterface[index].Key(), GetConversationIDBySessionType(friendsInfoOnLocalInterface[index].Key(), SingleChatType), SingleChatType)
		}
	}

	if len(sameA) > 0 {
		for _, index := range sameA {
			if friendInfoStruct, ok := friendsInfoOnServerInterface[index].Value().(friendInfo); ok {
				err = u.updateTheFriendInfo(friendInfoStruct.UID, friendInfoStruct.Name, friendInfoStruct.Comment, friendInfoStruct.Icon, friendInfoStruct.Gender, friendInfoStruct.Mobile, friendInfoStruct.Birth, friendInfoStruct.Email, friendInfoStruct.Ex)
				if err != nil {
					log(err.Error())
					return
				}
				jsonFriendInfo, _ := json.Marshal(friendInfoStruct)
				u.friendListener.OnFriendInfoChanged(string(jsonFriendInfo))
			}
		}
	}
}
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

func (u *UserRelated) getServerFriendList() ([]friendInfo, error) {
	resp, err := post2Api(getFriendListRouter, paramsCommonReq{OperationID: operationIDGenerator()}, u.token)
	if err != nil {
		return nil, err
	}
	var vgetFriendListResp getFriendListResp
	err = json.Unmarshal(resp, &vgetFriendListResp)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return nil, err
	}
	if vgetFriendListResp.ErrCode != 0 {
		sdkLog("errcode: ", vgetFriendListResp.ErrCode, "errmsg: ", vgetFriendListResp.ErrMsg)
		return nil, errors.New(vgetFriendListResp.ErrMsg)
	}
	return vgetFriendListResp.Data, nil
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
				log(err.Error())
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
				log(err.Error())
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

func (u *UserRelated) getServerFriendApplication() ([]applyUserInfo, error) {
	resp, err := post2Api(getFriendApplicationListRouter, paramsCommonReq{OperationID: operationIDGenerator()}, u.token)
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

func (u *UserRelated) doAcceptOrRefuseApplicationCall(sendUid string, flag int32) {
	sdkLog("doAcceptOrRefuseApplicationCall", sendUid, flag)

	var ui2GetUserInfo ui2ClientCommonReq
	ui2GetUserInfo.UidList = append(ui2GetUserInfo.UidList, sendUid)
	resp, err := post2Api(getUserInfoRouter, paramsGetUserInfo{UidList: ui2GetUserInfo.UidList, OperationID: operationIDGenerator()}, u.token)
	if err != nil {
		sdkLog("getUserInfo failed", err)
		return
	}
	var vgetUserInfoResp getUserInfoResp
	err = json.Unmarshal(resp, &vgetUserInfoResp)
	if err != nil {

	}
	if vgetUserInfoResp.ErrCode != 0 {
		sdkLog(vgetUserInfoResp.ErrCode, vgetUserInfoResp.ErrMsg)
		return
	}
	var appUserNode applyUserInfo
	appUserNode.Uid = vgetUserInfoResp.Data[0].Uid
	appUserNode.Name = vgetUserInfoResp.Data[0].Name
	appUserNode.Icon = vgetUserInfoResp.Data[0].Icon
	appUserNode.Gender = vgetUserInfoResp.Data[0].Gender
	appUserNode.Mobile = vgetUserInfoResp.Data[0].Mobile
	appUserNode.Birth = vgetUserInfoResp.Data[0].Birth
	appUserNode.Email = vgetUserInfoResp.Data[0].Email
	appUserNode.Ex = vgetUserInfoResp.Data[0].Ex
	appUserNode.Flag = flag

	jsonInfo, err := json.Marshal(appUserNode)
	if err != nil {
		sdkLog("doAcceptOrRefuseApplication json marshal failed")
		return
	}
	sdkLog(flag)
	if flag == 1 {
		u.friendListener.OnFriendApplicationListAccept(string(jsonInfo))
	}
	if flag == -1 {
		u.friendListener.OnFriendApplicationListReject(string(jsonInfo))
	}
}

func (u *UserRelated) refuseFriendApplication(uid2Refuse string) error {
	flag := -1
	resp, err := post2Api(addFriendResponse, paramsAddFriendResponse{Uid: uid2Refuse, OperationID: operationIDGenerator(), Flag: flag}, u.token)
	if err != nil {
		return err
	}
	var addFriendResp commonResp
	err = json.Unmarshal(resp, &addFriendResp)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return err
	}

	if addFriendResp.ErrCode != 0 {
		return errors.New(addFriendResp.ErrMsg)
	}

	u.syncFriendApplication()

	return nil
}

func (u *UserRelated) acceptFriendApplication(uid string) error {
	flag := 1
	resp, err := post2Api(addFriendResponse, paramsAddFriendResponse{Uid: uid, OperationID: operationIDGenerator(), Flag: flag}, u.token)
	if err != nil {
		return err
	}
	var addFriendResp commonResp
	err = json.Unmarshal(resp, &addFriendResp)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return err
	}
	if addFriendResp.ErrCode != 0 {
		sdkLog("errcode: ", addFriendResp.ErrCode, addFriendResp.ErrMsg)
		return errors.New(addFriendResp.ErrMsg)
	}

	u.syncFriendApplication()
	u.syncFriendList()
	n := NotificationContent{1, FriendAcceptTip, ""}
	u.autoSendMsg(u.createTextSystemMessage(n, AcceptFriendApplicationTip), uid, "", false, true, true)
	return nil
}

func (u *UserRelated) syncFriendApplication() {
	applicationListOnServer, err := u.getServerFriendApplication()
	if err != nil {
		return
	}
	applicationListOnServerInterface := make([]diff, 0)
	for _, v := range applicationListOnServer {
		applicationListOnServerInterface = append(applicationListOnServerInterface, v)
	}
	applicationListOnLocal, err := u.getLocalFriendApplication()
	if err != nil {
		return
	}
	applicationListOnLocalInterface := make([]diff, 0)
	for _, v := range applicationListOnLocal {
		applicationListOnLocalInterface = append(applicationListOnLocalInterface, v)
	}

	aInBNot, bInANot, sameA, _ := checkDiff(applicationListOnServerInterface, applicationListOnLocalInterface)

	if len(aInBNot) > 0 {
		for _, index := range aInBNot {
			if applicationListStruct, ok := applicationListOnServerInterface[index].Value().(applyUserInfo); ok {
				err = u.insertIntoTheUserToApplicationList(applicationListStruct)
				if err != nil {
					return
				}
			}
		}
	}

	if len(bInANot) > 0 {
		for _, index := range bInANot {
			err = u.delTheUserFromApplicationList(applicationListOnLocalInterface[index].Key())
			if err != nil {
				return
			}

		}
	}

	if len(sameA) > 0 {
		for _, index := range sameA {
			//interface--->struct
			if applicationListStruct, ok := applicationListOnServerInterface[index].Value().(applyUserInfo); ok {
				err = u.updateApplicationList(applicationListStruct)
				if err != nil {
					log(err.Error())
					return
				}
				jsonApplicationUserInfo, _ := json.Marshal(applicationListStruct)
				if applicationListStruct.Flag == 1 {
					_ = triggerCmdFriend()
					u.friendListener.OnFriendApplicationListAccept(string(jsonApplicationUserInfo))
				}
				if applicationListStruct.Flag == -1 {
					u.friendListener.OnFriendApplicationListReject(string(jsonApplicationUserInfo))
				}
			}
		}
	}
}

func (u *UserRelated) syncFriendList() {
	friendsInfoOnServer, err := u.getServerFriendList()
	if err != nil {
		return
	}
	friendsInfoOnServerInterface := make([]diff, 0)
	for _, v := range friendsInfoOnServer {
		friendsInfoOnServerInterface = append(friendsInfoOnServerInterface, v)
	}
	friendsInfoOnLocal, err := u.getLocalFriendList()
	if err != nil {
		return
	}
	friendsInfoOnLocalInterface := make([]diff, 0)
	for _, v := range friendsInfoOnLocal {
		friendsInfoOnLocalInterface = append(friendsInfoOnLocalInterface, v)
	}
	aInBNot, bInANot, sameA, _ := checkDiff(friendsInfoOnServerInterface, friendsInfoOnLocalInterface)
	if len(aInBNot) > 0 {
		for _, index := range aInBNot {
			if friendInfoStruct, ok := friendsInfoOnServerInterface[index].Value().(friendInfo); ok {
				sdkLog("insertIntoTheFriendToFriendInfo")
				err = u.insertIntoTheFriendToFriendInfo(friendInfoStruct.UID, friendInfoStruct.Name, friendInfoStruct.Comment, friendInfoStruct.Icon, friendInfoStruct.Gender, friendInfoStruct.Mobile, friendInfoStruct.Birth, friendInfoStruct.Email, friendInfoStruct.Ex)
				if err != nil {
					return
				}

			}
		}
	}

	if len(bInANot) > 0 {
		for _, index := range bInANot {
			sdkLog("delTheFriendFromFriendInfo")
			err = u.delTheFriendFromFriendInfo(friendsInfoOnLocalInterface[index].Key())
			if err != nil {
				log(err.Error())
				return
			}

		}
	}

	if len(sameA) > 0 {
		for _, index := range sameA {
			if friendInfoStruct, ok := friendsInfoOnServerInterface[index].Value().(friendInfo); ok {
				sdkLog("updateTheFriendInfo")
				err = u.updateTheFriendInfo(friendInfoStruct.UID, friendInfoStruct.Name, friendInfoStruct.Comment, friendInfoStruct.Icon, friendInfoStruct.Gender, friendInfoStruct.Mobile, friendInfoStruct.Birth, friendInfoStruct.Email, friendInfoStruct.Ex)
				if err != nil {
					log(err.Error())
					return
				}
			}
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
				log(err.Error())
				return
			}

		}
	}

	if len(bInANot) > 0 {
		for _, index := range bInANot {
			err = u.delTheUserFromBlackList(blackListOnLocalInterface[index].Key())
			if err != nil {
				log(err.Error())
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

func (u *UserRelated) addFriend(uiFriend paramsUiAddFriend) error {
	resp, err := post2Api(addFriendRouter, paramsAddFriend{UID: uiFriend.UID, OperationID: operationIDGenerator(), ReqMessage: uiFriend.ReqMessage}, u.token)
	if err != nil {
		return err
	}
	var cmResp commonResp
	err = json.Unmarshal(resp, &cmResp)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error(), resp)
		return err
	}
	if cmResp.ErrCode != 0 {
		return errors.New(cmResp.ErrMsg)
	}
	return nil
}
