package open_im_sdk

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

func (u *UserRelated) GetFriendList(callback Base, operationID string) {
	go func() {
		var filterLocalFriendList []LocalFriend
		localFriendList, err := u._getAllFriendList()
		checkErr(callback, err, operationID)
		localBlackUidList, err := u._getBlackListUid()
		checkErr(callback, err, operationID)
		for _, v := range localFriendList {
			if !isContain(v.FriendUserID, localBlackUidList) {
				filterLocalFriendList = append(filterLocalFriendList, v)
			}
		}
		callback.OnSuccess(structToJsonString(filterLocalFriendList))
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

func (u *UserRelated) AddToBlackList(callback Base, blackUid, operationID string) {
	go func() {
		u.addBlack(callback, blackUid, operationID)
		callback.OnSuccess("")
	}()
}

func (u *UserRelated) GetBlackList(callback Base, operationID string) {
	go func() {
		localBlackList, err := u._getBlackList()
		checkErr(callback, err, operationID)
		callback.OnSuccess(structToJsonString(localBlackList))
	}()
}

func (u *UserRelated) RemoveFromBlackList(callback Base, removeUid, operationID string) {
	go func() {
		u.removeBlack(callback, removeUid, operationID)
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
