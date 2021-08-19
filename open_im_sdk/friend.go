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

/*
func (fr *Friend) getCh() chan cmd2Value {
	return fr.ch
}
*/
/*
func (fr *Friend) work(c2v cmd2Value) {
	switch c2v.Cmd {
	case CmdFriend:
		if fr.friendListener == nil {
			sdkLog("friendListener is null")
			return
		}
		fr.doFriendList()
	case CmdBlackList:
		if fr.friendListener == nil {
			sdkLog("friendListener is null")
			return
		}
		fr.doBlackList()
	case CmdFriendApplication:

		if fr.friendListener == nil {
			sdkLog("friendListener is null")
			return
		}
		fr.doFriendApplication()
	case CmdAcceptFriend:
		if fr.friendListener == nil {
			sdkLog("friendListener is null")
			return
		}
		fr.doAcceptOrRefuseApplicationCall(c2v.Value.(string), 1)
	case CmdRefuseFriend:
		if fr.friendListener == nil {
			sdkLog("friendListener is null")
			return
		}
		sdkLog(c2v.Value.(string))
		fr.doAcceptOrRefuseApplicationCall(c2v.Value.(string), -1)
	case CmdAddFriend:
		if fr.friendListener == nil {
			sdkLog("friendListener is null")
			return
		}
		sdkLog(c2v.Value.(string))

	}
}
*/

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

func (fr *Friend) doFriendList() {
	friendsInfoOnServer, err := fr.getServerFriendList()
	if err != nil {
		return
	}
	friendsInfoOnServerInterface := make([]diff, 0)
	for _, v := range friendsInfoOnServer {
		friendsInfoOnServerInterface = append(friendsInfoOnServerInterface, v)
	}
	friendsInfoOnLocal, err := fr.getLocalFriendList()
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
				err = insertIntoTheFriendToFriendInfo(friendInfoStruct.UID, friendInfoStruct.Name, friendInfoStruct.Comment, friendInfoStruct.Icon, friendInfoStruct.Gender, friendInfoStruct.Mobile, friendInfoStruct.Birth, friendInfoStruct.Email, friendInfoStruct.Ex)
				if err != nil {
					log(err.Error())
					return
				}
				jsonFriendInfo, _ := json.Marshal(friendInfoStruct)
				fr.friendListener.OnFriendListAdded(string(jsonFriendInfo))
			}
		}
	}

	if len(bInANot) > 0 {
		for _, index := range bInANot {
			err = delTheFriendFromFriendInfo(friendsInfoOnLocalInterface[index].Key())
			if err != nil {
				log(err.Error())
				return
			}
			jsonFriendInfo, _ := json.Marshal(friendsInfoOnLocal[index])
			fr.friendListener.OnFriendListDeleted(string(jsonFriendInfo))
			_ = triggerCmdDeleteConversationAndMessage(friendsInfoOnLocalInterface[index].Key(), GetConversationIDBySessionType(friendsInfoOnLocalInterface[index].Key(), SingleChatType), SingleChatType)
		}
	}

	if len(sameA) > 0 {
		for _, index := range sameA {
			if friendInfoStruct, ok := friendsInfoOnServerInterface[index].Value().(friendInfo); ok {
				err = updateTheFriendInfo(friendInfoStruct.UID, friendInfoStruct.Name, friendInfoStruct.Comment, friendInfoStruct.Icon, friendInfoStruct.Gender, friendInfoStruct.Mobile, friendInfoStruct.Birth, friendInfoStruct.Email, friendInfoStruct.Ex)
				if err != nil {
					log(err.Error())
					return
				}
				jsonFriendInfo, _ := json.Marshal(friendInfoStruct)
				fr.friendListener.OnFriendInfoChanged(string(jsonFriendInfo))
			}
		}
	}
}
func (fr *Friend) getLocalFriendList() ([]friendInfo, error) {
	//Take out the friend list and judge whether it is in the blacklist again to prevent nested locks
	localFriendList, err := getLocalFriendList()
	if err != nil {
		return nil, err
	}
	for index, v := range localFriendList {
		//check friend is in blacklist
		blackUser, err := getBlackUsInfoByUid(v.UID)
		if err != nil {
			sdkLog(err.Error())
		}
		if blackUser.Uid != "" {
			localFriendList[index].IsInBlackList = 1
		}
	}
	return localFriendList, nil
}

func (fr *Friend) getServerFriendList() ([]friendInfo, error) {
	resp, err := post2Api(getFriendListRouter, paramsCommonReq{OperationID: operationIDGenerator()}, token)
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
func (fr *Friend) doBlackList() {

	blackListOnServer, err := fr.getServerBlackList()
	if err != nil {
		return
	}
	blackListOnServerInterface := make([]diff, 0)
	for _, blackUser := range blackListOnServer {
		blackListOnServerInterface = append(blackListOnServerInterface, blackUser)
	}

	blackListOnLocal, err := fr.getLocalBlackList()
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
			err = insertIntoTheUserToBlackList(blackListOnServer[index])
			if err != nil {
				log(err.Error())
				return
			}
			jsonAddBlackUserInfo, _ := json.Marshal(blackListOnServerInterface[index])
			fr.friendListener.OnBlackListAdd(string(jsonAddBlackUserInfo))
		}
	}

	if len(bInANot) > 0 {
		for _, index := range bInANot {
			err = delTheUserFromBlackList(blackListOnLocalInterface[index].Key())
			if err != nil {
				log(err.Error())
				return
			}
			jsonDelBlackUserInfo, _ := json.Marshal(blackListOnLocal[index])
			fr.friendListener.OnBlackListDeleted(string(jsonDelBlackUserInfo))
		}
	}
	if len(bInANot) > 0 || len(aInBNot) > 0 {
		_ = triggerCmdFriend()
	}

	if len(sameA) > 0 {
		for _, index := range sameA {
			//interface--->struct
			if blackListStruct, ok := blackListOnServerInterface[index].Value().(userInfo); ok {
				_ = updateBlackList(blackListStruct)
			}
		}
	}
}

func (fr *Friend) getServerBlackList() ([]userInfo, error) {
	resp, err := post2Api(getBlackListRouter, paramsCommonReq{OperationID: operationIDGenerator()}, token)
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

func (fr *Friend) getLocalBlackList() ([]userInfo, error) {
	return getLocalBlackList()
}

func (fr *Friend) doFriendApplication() {

	applicationListOnServer, err := fr.getServerFriendApplication()
	if err != nil {
		return
	}
	applicationListOnServerInterface := make([]diff, 0)
	for _, v := range applicationListOnServer {
		applicationListOnServerInterface = append(applicationListOnServerInterface, v)
	}
	applicationListOnLocal, err := fr.getLocalFriendApplication()
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
				err = insertIntoTheUserToApplicationList(applicationListStruct)
				if err != nil {
					log(err.Error())
					return
				}
				jsonAddApplicationUserInfo, _ := json.Marshal(applicationListStruct)
				sdkLog("OnFriendApplicationListAdded ")
				fr.friendListener.OnFriendApplicationListAdded(string(jsonAddApplicationUserInfo))
			}
		}
	}

	if len(bInANot) > 0 {
		for _, index := range bInANot {
			err = delTheUserFromApplicationList(applicationListOnLocalInterface[index].Key())
			if err != nil {
				log(err.Error())
				return
			}
			jsonDelApplicationUserInfo, _ := json.Marshal(applicationListOnLocal[index])
			fr.friendListener.OnFriendApplicationListDeleted(string(jsonDelApplicationUserInfo))
		}
	}

	if len(sameA) > 0 {
		for _, index := range sameA {
			//interface--->struct
			if applicationListStruct, ok := applicationListOnServerInterface[index].Value().(applyUserInfo); ok {
				err = updateApplicationList(applicationListStruct)
				if err != nil {
					log(err.Error())
					return
				}
				jsonApplicationUserInfo, _ := json.Marshal(applicationListStruct)
				if applicationListStruct.Flag == 1 {
					_ = triggerCmdFriend()
					fr.friendListener.OnFriendApplicationListAccept(string(jsonApplicationUserInfo))
				}
				if applicationListStruct.Flag == -1 {
					fr.friendListener.OnFriendApplicationListReject(string(jsonApplicationUserInfo))
				}
			}
		}
	}
}

func (fr *Friend) getServerFriendApplication() ([]applyUserInfo, error) {
	resp, err := post2Api(getFriendApplicationListRouter, paramsCommonReq{OperationID: operationIDGenerator()}, token)
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

func (fr *Friend) getServerSelfApplication() ([]applyUserInfo, error) {
	resp, err := post2Api(getSelfApplicationListRouter, paramsCommonReq{OperationID: operationIDGenerator()}, token)
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

func (fr *Friend) getLocalFriendApplication() ([]applyUserInfo, error) {
	return getLocalFriendApplication()
}

func (fr *Friend) doAcceptOrRefuseApplicationCall(sendUid string, flag int32) {
	sdkLog("doAcceptOrRefuseApplicationCall", sendUid, flag)

	var ui2GetUserInfo ui2ClientCommonReq
	ui2GetUserInfo.UidList = append(ui2GetUserInfo.UidList, sendUid)
	resp, err := post2Api(getUserInfoRouter, paramsGetUserInfo{UidList: ui2GetUserInfo.UidList, OperationID: operationIDGenerator()}, token)
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
		fr.friendListener.OnFriendApplicationListAccept(string(jsonInfo))
	}
	if flag == -1 {
		fr.friendListener.OnFriendApplicationListReject(string(jsonInfo))
	}
}

func (fr *Friend) refuseFriendApplication(uid2Refuse string) error {
	flag := -1
	resp, err := post2Api(addFriendResponse, paramsAddFriendResponse{Uid: uid2Refuse, OperationID: operationIDGenerator(), Flag: flag}, token)
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

	fr.syncFriendApplication()

	return nil
}

func (fr *Friend) acceptFriendApplication(uid string) error {
	flag := 1
	resp, err := post2Api(addFriendResponse, paramsAddFriendResponse{Uid: uid, OperationID: operationIDGenerator(), Flag: flag}, token)
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

	fr.syncFriendApplication()
	fr.syncFriendList()
	n := NotificationContent{1, FriendAcceptTip, ""}
	autoSendMsg(createTextSystemMessage(n, AcceptFriendApplicationTip), uid, "", false, true, true)
	return nil
}

func (fr *Friend) syncFriendApplication() {
	applicationListOnServer, err := fr.getServerFriendApplication()
	if err != nil {
		return
	}
	applicationListOnServerInterface := make([]diff, 0)
	for _, v := range applicationListOnServer {
		applicationListOnServerInterface = append(applicationListOnServerInterface, v)
	}
	applicationListOnLocal, err := fr.getLocalFriendApplication()
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
				err = insertIntoTheUserToApplicationList(applicationListStruct)
				if err != nil {
					return
				}
			}
		}
	}

	if len(bInANot) > 0 {
		for _, index := range bInANot {
			err = delTheUserFromApplicationList(applicationListOnLocalInterface[index].Key())
			if err != nil {
				return
			}

		}
	}

	if len(sameA) > 0 {
		for _, index := range sameA {
			//interface--->struct
			if applicationListStruct, ok := applicationListOnServerInterface[index].Value().(applyUserInfo); ok {
				err = updateApplicationList(applicationListStruct)
				if err != nil {
					log(err.Error())
					return
				}
				jsonApplicationUserInfo, _ := json.Marshal(applicationListStruct)
				if applicationListStruct.Flag == 1 {
					_ = triggerCmdFriend()
					fr.friendListener.OnFriendApplicationListAccept(string(jsonApplicationUserInfo))
				}
				if applicationListStruct.Flag == -1 {
					fr.friendListener.OnFriendApplicationListReject(string(jsonApplicationUserInfo))
				}
			}
		}
	}
}

func (fr *Friend) syncFriendList() {
	friendsInfoOnServer, err := fr.getServerFriendList()
	if err != nil {
		return
	}
	friendsInfoOnServerInterface := make([]diff, 0)
	for _, v := range friendsInfoOnServer {
		friendsInfoOnServerInterface = append(friendsInfoOnServerInterface, v)
	}
	friendsInfoOnLocal, err := fr.getLocalFriendList()
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
				err = insertIntoTheFriendToFriendInfo(friendInfoStruct.UID, friendInfoStruct.Name, friendInfoStruct.Comment, friendInfoStruct.Icon, friendInfoStruct.Gender, friendInfoStruct.Mobile, friendInfoStruct.Birth, friendInfoStruct.Email, friendInfoStruct.Ex)
				if err != nil {
					return
				}

			}
		}
	}

	if len(bInANot) > 0 {
		for _, index := range bInANot {
			sdkLog("delTheFriendFromFriendInfo")
			err = delTheFriendFromFriendInfo(friendsInfoOnLocalInterface[index].Key())
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
				err = updateTheFriendInfo(friendInfoStruct.UID, friendInfoStruct.Name, friendInfoStruct.Comment, friendInfoStruct.Icon, friendInfoStruct.Gender, friendInfoStruct.Mobile, friendInfoStruct.Birth, friendInfoStruct.Email, friendInfoStruct.Ex)
				if err != nil {
					log(err.Error())
					return
				}
			}
		}
	}
}

func (fr *Friend) syncBlackList() {

	blackListOnServer, err := fr.getServerBlackList()
	if err != nil {
		return
	}
	blackListOnServerInterface := make([]diff, 0)
	for _, blackUser := range blackListOnServer {
		blackListOnServerInterface = append(blackListOnServerInterface, blackUser)
	}

	blackListOnLocal, err := fr.getLocalBlackList()
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
			err = insertIntoTheUserToBlackList(blackListOnServer[index])
			if err != nil {
				log(err.Error())
				return
			}

		}
	}

	if len(bInANot) > 0 {
		for _, index := range bInANot {
			err = delTheUserFromBlackList(blackListOnLocalInterface[index].Key())
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
				_ = updateBlackList(blackListStruct)
			}
		}
	}
}

func (fr *Friend) addFriend(uiFriend paramsUiAddFriend) error {
	resp, err := post2Api(addFriendRouter, paramsAddFriend{UID: uiFriend.UID, OperationID: operationIDGenerator(), ReqMessage: uiFriend.ReqMessage}, token)
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
