package open_im_sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (u *UserRelated) GetFriendsInfo(callback Base, uidList string) {
	if callback == nil || uidList == "" {
		sdkLog("uidList or callback is nil")
		return
	}
	go func() {
		fList, err := u.getLocalFriendList()
		if err != nil {
			sdkLog("getLocalFriendList failed, ", err.Error())
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}
		var sList []string
		e := json.Unmarshal([]byte(uidList), &sList)
		if e != nil {
			callback.OnError(ErrCodeFriend, e.Error())
			sdkLog("Unmarshal failed, ", e.Error())
			return
		}
		mapFriend := make(map[string]friendInfo)
		for _, v := range fList {
			mapFriend[v.UID] = v
		}

		result := make([]friendInfo, 0)
		for _, v := range sList {
			k, ok := mapFriend[v]
			if ok {
				result = append(result, k)
			}
		}

		sr, ee := json.Marshal(result)
		if ee != nil {
			callback.OnError(ErrCodeFriend, ee.Error())
			sdkLog("Marshal failed, ", ee.Error())
			return
		}
		callback.OnSuccess(string(sr))
	}()
}

func (u *UserRelated) AddFriend(callback Base, paramsReq string) {
	go func() {
		var uiFriend paramsUiAddFriend
		err := json.Unmarshal([]byte(paramsReq), &uiFriend)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			sdkLog("unmarshal failed, ", err.Error(), paramsReq)
			return
		}
		err = u.addFriend(uiFriend)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}
		callback.OnSuccess("")
	}()
}

func (u *UserRelated) GetFriendApplicationList(callback Base) {
	go func() {
		list, err := u.getLocalFriendApplication()
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			sdkLog("getLocalFriendApplication failed ", err.Error())
			return
		}
		slist, err := json.Marshal(list)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			sdkLog("Marshal failed ", err.Error())
			return
		}
		callback.OnSuccess(string(slist))
	}()
}

func (u *UserRelated) AcceptFriendApplication(callback Base, uid string) {
	//FriendApplication(callback, info, 1)
	go func() {
		var uid2Accept string
		err := json.Unmarshal([]byte(uid), &uid2Accept)
		if err != nil {
			sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}

		err = u.acceptFriendApplication(uid2Accept)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}

		fInfo, err := u.getFriendInfoByFriendUid(uid2Accept)

		blackUser, err := u.getBlackUsInfoByUid(uid2Accept)
		if err != nil {
			sdkLog(err.Error())
		}
		if blackUser.Uid != "" {
			fInfo.IsInBlackList = 1
		}
		if err == nil && fInfo.UID != "" {
			jsonInfo, err := json.Marshal(fInfo)
			if err == nil {
				u.friendListener.OnFriendListAdded(string(jsonInfo))
			}
		}
		callback.OnSuccess("")
	}()
}

func (u *UserRelated) RefuseFriendApplication(callback Base, uid string) {
	//	FriendApplication(callback, uid, -1)
	go func() {
		var uid2Refuse string
		err := json.Unmarshal([]byte(uid), &uid2Refuse)
		if err != nil {
			sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}

		err = u.refuseFriendApplication(uid2Refuse)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}
		callback.OnSuccess("")
	}()
}

func (u *UserRelated) CheckFriend(callback Base, uidList string) {
	go func() {
		fList, err := u.getLocalFriendList()
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			log(fmt.Sprintf("CheckFriend ErrCodeFriend err = %s", err.Error()))
			return
		}
		var ui2UidList []string
		e := json.Unmarshal([]byte(uidList), &ui2UidList)
		if e != nil {
			callback.OnError(ErrCodeFriend, e.Error())

			log(fmt.Sprintf("CheckFriend Unmarshal err = %s", e.Error()))

			return
		}
		mapFriend := make(map[string]int32)
		for _, v := range fList {
			mapFriend[v.UID] = 1
		}

		result := make([]Uid2Flag, 0)
		for _, v := range ui2UidList {
			result = append(result, Uid2Flag{Uid: v, Flag: mapFriend[v]})
		}

		sr, ee := json.Marshal(result)
		if ee != nil {
			callback.OnError(ErrCodeFriend, ee.Error())
			log(fmt.Sprintf("CheckFriend Marshal err = %s", ee.Error()))
			return
		}
		callback.OnSuccess(string(sr))
	}()
}

func (ur *UserRelated) DeleteFromFriendList(deleteUid string, callback Base) {
	go func() {
		var dUid string
		er := json.Unmarshal([]byte(deleteUid), &dUid)
		if er != nil {
			callback.OnError(ErrCodeFriend, er.Error())
			sdkLog("Unmarshal failed, ", er.Error(), deleteUid)
			return
		}

		resp, err := post2Api(deleteFriendRouter, paramsDeleteFriend{Uid: dUid, OperationID: operationIDGenerator()}, ur.token)
		if err != nil {
			callback.OnError(http.StatusInternalServerError, err.Error())
			sdkLog("post2Api failed, ", err.Error())
			return
		}
		var deleteFriendResp commonResp
		_ = json.Unmarshal(resp, &deleteFriendResp)
		if deleteFriendResp.ErrCode != 0 {
			callback.OnError(deleteFriendResp.ErrCode, deleteFriendResp.ErrMsg)

			log(fmt.Sprintf("DeleteFromFriendList Unmarshal errcode = %d", deleteFriendResp.ErrCode))

			return
		}
		//_ = triggerCmdFriend()

		ur.syncFriendList()
		u, err := ur.getUserInfoByUid(dUid)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			log(fmt.Sprintf("getUserInfoByUid  err = %s", err.Error()))
			return
		}
		f := friendInfo{
			UID:    u.Uid,
			Name:   u.Name,
			Icon:   u.Icon,
			Gender: u.Gender,
			Mobile: u.Mobile,
			Birth:  u.Birth,
			Email:  u.Email,
			Ex:     u.Ex,
		}
		ur.friendListener.OnFriendListDeleted(structToJsonString(f))
		//_ = ur.triggerCmdDeleteConversationAndMessage(dUid, GetConversationIDBySessionType(dUid, SingleChatType), SingleChatType)
		callback.OnSuccess("")
	}()
}

func (u *UserRelated) GetFriendList(callback Base) {
	go func() {
		list, err := u.getLocalFriendList()
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
		} else {
			jlist, e := json.Marshal(list)
			if e != nil {
				callback.OnError(ErrCodeFriend, e.Error())
			} else {
				callback.OnSuccess(string(jlist))
			}
		}
	}()
}

func (u *UserRelated) SetFriendInfo(comment string, callback Base) {
	go func() {
		var uid2comm uid2Comment
		er := json.Unmarshal([]byte(comment), &uid2comm)
		if er != nil {
			callback.OnError(ErrCodeFriend, er.Error())

			log(fmt.Sprintf("SetFriendInfo ErrCodeFriend err = %s", er.Error()))

			return
		}
		resp, err := post2Api(setFriendComment, paramsSetFriendInfo{Uid: uid2comm.Uid, OperationID: operationIDGenerator(), Comment: uid2comm.Comment}, u.token)
		if err != nil {
			callback.OnError(http.StatusInternalServerError, err.Error())

			log(fmt.Sprintf("SetFriendInfo StatusInternalServerError err = %s", er.Error()))

			return
		}
		var cResp commonResp
		err = json.Unmarshal(resp, &cResp)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())

			log(fmt.Sprintf("SetFriendInfo Unmarshal err = %s", er.Error()))

			return
		}
		if cResp.ErrCode != 0 {
			callback.OnError(ErrCodeFriend, cResp.ErrMsg)

			return
		}
		u.syncFriendList()
		callback.OnSuccess("")
		c := ConversationStruct{
			ConversationID: GetConversationIDBySessionType(uid2comm.Uid, SingleChatType),
		}
		faceUrl, name, err := u.getUserNameAndFaceUrlByUid(uid2comm.Uid)
		if err != nil {
			sdkLog("getUserNameAndFaceUrlByUid err:", err)
			return
		}
		c.FaceURL = faceUrl
		c.ShowName = name
		u.doUpdateConversation(cmd2Value{Value: updateConNode{c.ConversationID, UpdateFaceUrlAndNickName, c}})
		u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{c.ConversationID}}})

		//FriendObj.friendListener.OnFriendInfoChanged(structToJsonString(friendResp.Data))
		//_ = triggerCmdFriend()
	}()
}

func (u *UserRelated) AddToBlackList(callback Base, blackUid string) {
	go func() {
		var uid string
		er := json.Unmarshal([]byte(blackUid), &uid)
		if er != nil {
			callback.OnError(ErrCodeFriend, er.Error())

			log(fmt.Sprintf("AddToBlackList ErrCodeFriend err = %s", er.Error()))

			return
		}
		resp, err := post2Api(addBlackListRouter, paramsAddBlackList{UID: uid, OperationID: operationIDGenerator()}, u.token)
		if err != nil {
			callback.OnError(http.StatusInternalServerError, err.Error())
			log(fmt.Sprintf("AddToBlackList StatusInternalServerError err = %s", er.Error()))
			return
		}

		var addToBlackListResp commonResp
		err = json.Unmarshal(resp, &addToBlackListResp)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())

			log(fmt.Sprintf("AddToBlackList ErrCodeFriend err = %s", err.Error()))

			return
		}

		if addToBlackListResp.ErrCode != 0 {
			callback.OnError(addToBlackListResp.ErrCode, addToBlackListResp.ErrMsg)
		}

		user, err := u.getUserInfoByUid(uid)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			log(fmt.Sprintf("getUserInfoByUid err = %s", err.Error()))
			return
		}

		u.syncBlackList()

		bUser, err := json.Marshal(user)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())

			log(fmt.Sprintf("AddToBlackList ErrCodeFriend err = %s", err.Error()))

			return
		}

		u.friendListener.OnBlackListAdd(string(bUser))
		callback.OnSuccess("")
	}()
}

func (u *UserRelated) GetBlackList(callback Base) {
	go func() {
		list, err := u.getServerBlackList()
		if err == nil {
			jlist, e := json.Marshal(list)
			if e == nil {
				callback.OnSuccess(string(jlist))
				return
			}
		}

		list, err = u.getLocalBlackList()
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			log(fmt.Sprintf("getBlackList ErrCodeFriend err = %s", err.Error()))
		} else {
			jlist, e := json.Marshal(list)
			if e != nil {
				callback.OnError(ErrCodeFriend, e.Error())
			} else {
				callback.OnSuccess(string(jlist))
			}
		}
	}()
}

func (u *UserRelated) DeleteFromBlackList(callback Base, deleteUid string) {
	go func() {
		var duid string
		er := json.Unmarshal([]byte(deleteUid), &duid)
		if er != nil {
			callback.OnError(ErrCodeFriend, er.Error())

			log(fmt.Sprintf("DeleteFromBlackList ErrCodeFriend err = %s", er.Error()))

			return
		}
		resp, err := post2Api(removeBlackListRouter, paramsRemoveBlackList{UID: duid, OperationID: operationIDGenerator()}, u.token)
		if err != nil {
			callback.OnError(http.StatusInternalServerError, err.Error())
			log(fmt.Sprintf("DeleteFromBlackList StatusInternalServerError err = %s", err.Error()))
			return
		}
		var removeToBlackList commonResp
		err = json.Unmarshal(resp, &removeToBlackList)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			log(fmt.Sprintf("DeleteFromBlackList ErrCodeFriend err = %s", err.Error()))
			return
		}
		if removeToBlackList.ErrCode != 0 {
			callback.OnError(removeToBlackList.ErrCode, removeToBlackList.ErrMsg)
		}

		user, err := u.getUserInfoByUid(duid)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			log(fmt.Sprintf("getUserInfoByUid err = %s", err.Error()))
			return
		}

		u.syncBlackList()

		bUser, err := json.Marshal(user)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())

			log(fmt.Sprintf("AddToBlackList ErrCodeFriend err = %s", err.Error()))

			return
		}

		u.friendListener.OnBlackListDeleted(string(bUser))
		callback.OnSuccess("")
	}()
}

type OnFriendshipListener interface {
	OnFriendApplicationListAdded(applyUserInfo string)
	OnFriendApplicationListDeleted(applyUserInfo string)
	OnFriendApplicationListAccept(applyUserInfo string)
	OnFriendApplicationListReject(applyUserInfo string)
	OnFriendListAdded(friendInfo string)
	OnFriendListDeleted(friendInfo string)
	OnBlackListAdd(userInfo string)
	OnBlackListDeleted(userInfo string)
	OnFriendInfoChanged(friendInfo string)
}

func (u *UserRelated) SetFriendListener(listener OnFriendshipListener) bool {
	if listener == nil {
		return false
	}
	u.friendListener = listener
	return true
}

func (u *UserRelated) ForceSyncFriendApplication() {
	u.syncFriendApplication()
}

func (u *UserRelated) ForceSyncSelfFriendApplication() {
	u.syncSelfFriendApplication()
}

func (u *UserRelated) ForceSyncFriend() {
	u.syncFriendList()
}

func (u *UserRelated) ForceSyncBlackList() {
	u.syncBlackList()
}
