package open_im_sdk

func (u *UserRelated) GetDesignatedFriendsInfo(callback Base, params string, operationID string) {
	NewInfo(operationID, GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		NewInfo(operationID, "GetDesignatedFriendsInfo args: ", params)
		var unmarshalList GetDesignatedFriendsInfoParams
		u.jsonUnmarshal(params, &unmarshalList, callback, operationID)
		result := u.getDesignatedFriendsInfo(callback, unmarshalList, operationID)
		callback.OnSuccess(structToJsonStringDefault(result))
		NewInfo(operationID, "GetDesignatedFriendsInfo callback: ", structToJsonString(result), result)
	}()
}

func (u *UserRelated) AddFriend(callback Base, params string, operationID string) {
	NewInfo(operationID, GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		NewInfo(operationID, "AddFriend args: ", params)
		var unmarshalAddFriendParams AddFriendParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalAddFriendParams, callback, operationID)
		u.addFriend(callback, unmarshalAddFriendParams, operationID)
		callback.OnSuccess(structToJsonString(AddFriendCallback))
		NewInfo(operationID, "AddFriend callback: ", structToJsonString(AddFriendCallback))
	}()
}

func (u *UserRelated) GetRecvFriendApplicationList(callback Base, operationID string) {
	NewInfo(operationID, GetSelfFuncName())
	if callback == nil {
		return
	}
	go func() {
		NewInfo(operationID, "GetRecvFriendApplicationList args: ")
		result := u.getRecvFriendApplicationList(callback, operationID)
		callback.OnSuccess(structToJsonString(result))
		NewInfo(operationID, "GetRecvFriendApplicationList callback: ", structToJsonString(result))
	}()
}

func (u *UserRelated) GetSendFriendApplicationList(callback Base, operationID string) {
	NewInfo(operationID, GetSelfFuncName())
	if callback == nil {
		return
	}
	go func() {
		NewInfo(operationID, "GetSendFriendApplicationList args: ")
		result := u.getSendFriendApplicationList(callback, operationID)
		callback.OnSuccess(structToJsonString(result))
		NewInfo(operationID, "GetSendFriendApplicationList callback: ", structToJsonString(result))
	}()
}

func (u *UserRelated) AcceptFriendApplication(callback Base, params string, operationID string) {
	NewInfo(operationID, GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		NewInfo(operationID, "AcceptFriendApplication args: ", params)
		var unmarshalParams ProcessFriendApplicationParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback, operationID)
		u.processFriendApplication(callback, unmarshalParams, 1, operationID)
		callback.OnSuccess(structToJsonString(ProcessFriendApplicationCallback))
		NewInfo(operationID, "AcceptFriendApplication callback: ", structToJsonString(ProcessFriendApplicationCallback))
	}()
}

func (u *UserRelated) RefuseFriendApplication(callback Base, params string, operationID string) {
	NewInfo(operationID, GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		NewInfo(operationID, "RefuseFriendApplication args: ", params)
		var unmarshalParams ProcessFriendApplicationParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback, operationID)
		u.processFriendApplication(callback, unmarshalParams, -1, operationID)
		callback.OnSuccess(structToJsonString(ProcessFriendApplicationCallback))
		NewInfo(operationID, "RefuseFriendApplication callback: ", structToJsonString(ProcessFriendApplicationCallback))
	}()
}

func (u *UserRelated) CheckFriend(callback Base, params string, operationID string) {
	NewInfo(operationID, GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		NewInfo(operationID, "CheckFriend args: ", params)
		var unmarshalParams CheckFriendParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback, operationID)
		result := u.checkFriend(callback, unmarshalParams, operationID)
		callback.OnSuccess(structToJsonString(result))
		NewInfo(operationID, "CheckFriend callback: ", structToJsonString(result))
	}()
}

func (u *UserRelated) DeleteFriend(callback Base, friendUserID string, operationID string) {
	NewInfo(operationID, GetSelfFuncName(), friendUserID)
	if callback == nil {
		return
	}
	go func() {
		NewInfo(operationID, "DeleteFriend args: ", friendUserID)
		u.deleteFriend(friendUserID, callback, operationID)
		callback.OnSuccess(structToJsonString(DeleteFriendCallback))
		NewInfo(operationID, "DeleteFriend callback: ", structToJsonString(DeleteFriendCallback))
	}()
}

func (u *UserRelated) GetFriendList(callback Base, operationID string) {
	NewInfo(operationID, GetSelfFuncName())
	if callback == nil {
		return
	}
	go func() {
		NewInfo(operationID, "GetFriendList args: ")
		var filterLocalFriendList []*LocalFriend
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
		NewInfo(operationID, "GetFriendList callback: ", structToJsonString(filterLocalFriendList))
	}()
}

func (u *UserRelated) SetFriendRemark(callback Base, params string, operationID string) {
	NewInfo(operationID, GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		NewInfo(operationID, "SetFriendRemark args: ", params)
		var unmarshalParams SetFriendRemarkParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback, operationID)
		u.setFriendRemark(unmarshalParams, callback, operationID)
		callback.OnSuccess(structToJsonString(SetFriendRemarkCallback))
		NewInfo(operationID, "SetFriendRemark callback: ", structToJsonString(SetFriendRemarkCallback))
	}()
}

func (u *UserRelated) AddBlack(callback Base, blackUserID, operationID string) {
	NewInfo(operationID, GetSelfFuncName(), blackUserID)
	go func() {
		NewInfo(operationID, "AddToBlackList args: ", blackUserID)
		u.addBlack(callback, blackUserID, operationID)
		callback.OnSuccess("")
		NewInfo(operationID, "AddToBlackList callback: ")
	}()
}

func (u *UserRelated) GetBlackList(callback Base, operationID string) {
	NewInfo(operationID, GetSelfFuncName())
	go func() {
		NewInfo(operationID, "GetBlackList args: ")
		localBlackList, err := u._getBlackList()
		checkErr(callback, err, operationID)
		callback.OnSuccess(structToJsonString(localBlackList))
		NewInfo(operationID, "GetBlackList callback: ", structToJsonString(localBlackList))
	}()
}

func (u *UserRelated) RemoveBlack(callback Base, removeUserID, operationID string) {
	NewInfo(operationID, GetSelfFuncName(), removeUserID)
	go func() {
		NewInfo(operationID, "RemoveBlack args: ", removeUserID)
		u.removeBlack(callback, removeUserID, operationID)
		callback.OnSuccess("")
		NewInfo(operationID, "RemoveBlack callback")
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
