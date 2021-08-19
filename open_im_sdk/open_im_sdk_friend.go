package open_im_sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetFriendsInfo(callback Base, uidList string) {
	if callback == nil {
		return
	}
	go func() {
		fList, err := FriendObj.getLocalFriendList()
		if err != nil {
			sdkLog("getLocalFriendList failed, ", err.Error())
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}
		var sList []string
		e := json.Unmarshal([]byte(uidList), &sList)
		if e != nil {
			callback.OnError(ErrCodeFriend, e.Error())
			log(fmt.Sprintf("GetFriendsInfo err = %s", e.Error()))
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
			log(fmt.Sprintf("GetFriendsInfo err = %s", ee.Error()))
			return
		}
		callback.OnSuccess(string(sr))
	}()
}

func AddFriend(callback Base, paramsReq string) {
	go func() {
		var uiFriend paramsUiAddFriend
		err := json.Unmarshal([]byte(paramsReq), &uiFriend)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			sdkLog("unmarshal failed, ", err.Error(), paramsReq)
			return
		}
		err = FriendObj.addFriend(uiFriend)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}
		callback.OnSuccess("")
	}()
}

func GetFriendApplicationList(callback Base) {
	go func() {
		list, err := FriendObj.getLocalFriendApplication()
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			log(fmt.Sprintf("GetFriendApplicationList ErrCodeFriend err = %s", err.Error()))
			return
		}
		slist, err := json.Marshal(list)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			log(fmt.Sprintf("GetFriendApplicationList ErrCodeFriend err = %s", err.Error()))
			return
		}
		callback.OnSuccess(string(slist))
	}()
}

func AcceptFriendApplication(callback Base, uid string) {
	//FriendApplication(callback, info, 1)
	go func() {
		var uid2Accept string
		err := json.Unmarshal([]byte(uid), &uid2Accept)
		if err != nil {
			sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}

		err = FriendObj.acceptFriendApplication(uid2Accept)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}

		fInfo, err := getFriendInfoByFriendUid(uid2Accept)

		blackUser, err := getBlackUsInfoByUid(uid2Accept)
		if err != nil {
			sdkLog(err.Error())
		}
		if blackUser.Uid != "" {
			fInfo.IsInBlackList = 1
		}
		if err == nil && fInfo.UID != "" {
			jsonInfo, err := json.Marshal(fInfo)
			if err == nil {
				FriendObj.friendListener.OnFriendListAdded(string(jsonInfo))
			}
		}
		callback.OnSuccess("")
	}()
}

func RefuseFriendApplication(callback Base, uid string) {
	//	FriendApplication(callback, uid, -1)
	go func() {
		var uid2Refuse string
		err := json.Unmarshal([]byte(uid), &uid2Refuse)
		if err != nil {
			sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}

		err = FriendObj.refuseFriendApplication(uid2Refuse)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}
		callback.OnSuccess("")
	}()
}

func CheckFriend(callback Base, uidList string) {
	go func() {
		fList, err := FriendObj.getLocalFriendList()
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

func DeleteFromFriendList(deleteUid string, callback Base) {
	go func() {
		var dUid string
		er := json.Unmarshal([]byte(deleteUid), &dUid)
		if er != nil {
			callback.OnError(ErrCodeFriend, er.Error())

			log(fmt.Sprintf("DeleteFromFriendList Unmarshal err = %s", er.Error()))

			return
		}

		resp, err := post2Api(deleteFriendRouter, paramsDeleteFriend{Uid: dUid, OperationID: operationIDGenerator()}, token)
		if err != nil {
			callback.OnError(http.StatusInternalServerError, err.Error())

			log(fmt.Sprintf("DeleteFromFriendList StatusInternalServerError err = %s", er.Error()))

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

		FriendObj.syncFriendList()
		u, err := getUserInfoByUid(dUid)
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
		FriendObj.friendListener.OnFriendListDeleted(structToJsonString(f))
		_ = triggerCmdDeleteConversationAndMessage(dUid, GetConversationIDBySessionType(dUid, SingleChatType), SingleChatType)
		callback.OnSuccess("")
	}()
}

func FriendApplication(callback Base, uid string, flag int) {
	go func() {
		var uid2Accept ui2AcceptFriend
		err := json.Unmarshal([]byte(uid), &uid2Accept)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			return
		}
		resp, err := post2Api(addFriendResponse, paramsAddFriendResponse{Uid: uid2Accept.UID, OperationID: operationIDGenerator(), Flag: flag}, token)
		if err != nil {
			callback.OnError(http.StatusInternalServerError, err.Error())

			return
		}

		var addFriendResp commonResp
		err = json.Unmarshal(resp, &addFriendResp)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			log(fmt.Sprintf("FriendApplication ErrCodeFriend err = %d", err.Error()))
			return
		}

		if addFriendResp.ErrCode == 0 {
			callback.OnSuccess("")
			_ = triggerCmdFriendApplication()
			if flag == 1 {
				_ = triggerCmdFriend()
				//Trigger the latest message in the conversation and display it on the conversation box
				//m := &MsgStruct{
				//	Content: "You are already friends, come and chat together",
				//}
				//_ = triggerCmdUpdateConversation(updateConNode{ConId: getSingleConversationID(uid2Accept.UID), Action: 6, Args: ConversationStruct{
				//	ConversationID:   getSingleConversationID(uid2Accept.UID),
				//	ConversationType: SingleChatType,
				//	UserID:           uid2Accept.UID,
				//	RecvMsgOpt:       1,
				//	LatestMsg:        structToJsonString(m),
				//}})
			}
		} else {
			callback.OnError(addFriendResp.ErrCode, addFriendResp.ErrMsg)
		}
	}()
}

func GetFriendList(callback Base) {
	go func() {
		list, err := FriendObj.getLocalFriendList()
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

func SetFriendInfo(comment string, callback Base) {
	go func() {
		var uid2comm uid2Comment
		er := json.Unmarshal([]byte(comment), &uid2comm)
		if er != nil {
			callback.OnError(ErrCodeFriend, er.Error())

			log(fmt.Sprintf("SetFriendInfo ErrCodeFriend err = %s", er.Error()))

			return
		}
		resp, err := post2Api(setFriendComment, paramsSetFriendInfo{Uid: uid2comm.Uid, OperationID: operationIDGenerator(), Comment: uid2comm.Comment}, token)
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
		FriendObj.syncFriendList()
		callback.OnSuccess("")
		//FriendObj.friendListener.OnFriendInfoChanged(structToJsonString(friendResp.Data))
		//_ = triggerCmdFriend()
	}()
}

func AddToBlackList(callback Base, blackUid string) {
	go func() {
		var uid string
		er := json.Unmarshal([]byte(blackUid), &uid)
		if er != nil {
			callback.OnError(ErrCodeFriend, er.Error())

			log(fmt.Sprintf("AddToBlackList ErrCodeFriend err = %s", er.Error()))

			return
		}
		resp, err := post2Api(addBlackListRouter, paramsAddBlackList{UID: uid, OperationID: operationIDGenerator()}, token)
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

		user, err := getUserInfoByUid(uid)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			log(fmt.Sprintf("getUserInfoByUid err = %s", err.Error()))
			return
		}

		FriendObj.syncBlackList()

		bUser, err := json.Marshal(user)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())

			log(fmt.Sprintf("AddToBlackList ErrCodeFriend err = %s", err.Error()))

			return
		}

		FriendObj.friendListener.OnBlackListAdd(string(bUser))
		callback.OnSuccess("")
	}()
}

func GetBlackList(callback Base) {
	go func() {
		list, err := FriendObj.getServerBlackList()
		if err == nil {
			jlist, e := json.Marshal(list)
			if e == nil {
				callback.OnSuccess(string(jlist))
				return
			}
		}

		list, err = FriendObj.getLocalBlackList()
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

func DeleteFromBlackList(callback Base, deleteUid string) {
	go func() {
		var duid string
		er := json.Unmarshal([]byte(deleteUid), &duid)
		if er != nil {
			callback.OnError(ErrCodeFriend, er.Error())

			log(fmt.Sprintf("DeleteFromBlackList ErrCodeFriend err = %s", er.Error()))

			return
		}
		resp, err := post2Api(removeBlackListRouter, paramsRemoveBlackList{UID: duid, OperationID: operationIDGenerator()}, token)
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

		user, err := getUserInfoByUid(duid)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())
			log(fmt.Sprintf("getUserInfoByUid err = %s", err.Error()))
			return
		}

		FriendObj.syncBlackList()

		bUser, err := json.Marshal(user)
		if err != nil {
			callback.OnError(ErrCodeFriend, err.Error())

			log(fmt.Sprintf("AddToBlackList ErrCodeFriend err = %s", err.Error()))

			return
		}

		FriendObj.friendListener.OnBlackListDeleted(string(bUser))
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

func SetFriendListener(listener OnFriendshipListener) bool {
	if listener == nil {
		return false
	}
	FriendObj.friendListener = listener
	return true
}

func ForceSyncFriendApplication() {
	FriendObj.syncFriendApplication()
}

func ForceSyncFriend() {
	FriendObj.syncFriendList()
}

func ForceSyncBlackList() {
	FriendObj.syncBlackList()
}
