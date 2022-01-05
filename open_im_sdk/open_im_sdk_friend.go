package open_im_sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (u *UserRelated) GetFriendsInfo(callback Base, friendUserIDList string, operationID string) {
	if callback == nil || friendUserIDList == "" {
		sdkLog("uidList or callback is nil")
		return
	}
	go func() {
		var unmarshalList GetDesignatedFriendsInfoParams
		u.jsonUnmarshalAndArgsValidate(friendUserIDList, &unmarshalList, callback)
		result := u.getDesignatedFriendsInfo(callback, unmarshalList, operationID)
		callback.OnSuccess(structToJsonString(result))
	}()
}

func (u *UserRelated) AddFriend(callback Base, paramsReq string, operationID string) {
	go func() {
		var unmarshalAddFriendParams AddFriendParams
		u.jsonUnmarshalAndArgsValidate(paramsReq, &unmarshalAddFriendParams, callback)
		u.addFriend(callback, unmarshalAddFriendParams, operationID)
		callback.OnSuccess(structToJsonString(AddFriendCallback{}))
	}()
}

func (u *UserRelated) GetRecvFriendApplicationList(callback Base, operationID string) {
	go func() {
		u.getRecvFriendApplicationList(callback, operationID)
	}()
}

func (u *UserRelated) GetSendFriendApplicationList(callback Base, operationID string) {
	go func() {
		u.getSendFriendApplicationList(callback, operationID)
	}()
}

func (u *UserRelated) AcceptFriendApplication(callback Base, params string, operationID string) {
	go func() {
		var unmarshalParams ProcessFriendApplicationParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback)
		u.processFriendApplication(callback, unmarshalParams, 1, operationID)
		callback.OnSuccess(structToJsonString(ProcessFriendApplicationCallback{}))
	}()
}

func (u *UserRelated) RefuseFriendApplication(callback Base, params string, operationID string) {
	go func() {
		var unmarshalParams ProcessFriendApplicationParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback)
		u.processFriendApplication(callback, unmarshalParams, -1, operationID)
		callback.OnSuccess(structToJsonString(ProcessFriendApplicationCallback{}))
	}()
}

func (u *UserRelated) CheckFriend(callback Base, params string, operationID string) {
	go func() {
		var unmarshalParams CheckFriendParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback)
		result := u.checkFriend(callback, unmarshalParams, operationID)
		callback.OnSuccess(structToJsonString(result))
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
