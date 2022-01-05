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

func (u *UserRelated) DeleteFromFriendList(callback Base, friendUserID string, operationID string) {
	go func() {
		u.deleteFriend(friendUserID, callback, operationID)
		callback.OnSuccess(structToJsonString(DeleteFriendCallback{}))
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

func (u *UserRelated) SetFriendRemark(params string, callback Base, operationID string) {
	go func() {
		var unmarshalParams SetFriendRemarkParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback)
		u.setFriendRemark(unmarshalParams, callback, operationID)
		callback.OnSuccess(structToJsonString(SetFriendRemarkCallback{}))
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
