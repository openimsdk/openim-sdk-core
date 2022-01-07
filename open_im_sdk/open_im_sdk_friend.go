package open_im_sdk

import (
	"open_im_sdk/open_im_sdk/log"
	"open_im_sdk/open_im_sdk/sdk_params_callback"
	"open_im_sdk/open_im_sdk/utils"
)

func (u *UserRelated) GetDesignatedFriendsInfo(callback Base, params string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetDesignatedFriendsInfo args: ", params)
		var unmarshalList sdk_params_callback.GetDesignatedFriendsInfoParams
		u.jsonUnmarshal(params, &unmarshalList, callback, operationID)
		result := u.getDesignatedFriendsInfo(callback, unmarshalList, operationID)
		callback.OnSuccess(utils.structToJsonStringDefault(result))
		log.NewInfo(operationID, "GetDesignatedFriendsInfo callback: ", utils.structToJsonString(result), result)
	}()
}

func (u *UserRelated) AddFriend(callback Base, params string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "AddFriend args: ", params)
		var unmarshalAddFriendParams sdk_params_callback.AddFriendParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalAddFriendParams, callback, operationID)
		u.addFriend(callback, unmarshalAddFriendParams, operationID)
		callback.OnSuccess(utils.structToJsonString(sdk_params_callback.AddFriendCallback))
		log.NewInfo(operationID, "AddFriend callback: ", utils.structToJsonString(sdk_params_callback.AddFriendCallback))
	}()
}

func (u *UserRelated) GetRecvFriendApplicationList(callback Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetRecvFriendApplicationList args: ")
		result := u.getRecvFriendApplicationList(callback, operationID)
		callback.OnSuccess(utils.structToJsonString(result))
		log.NewInfo(operationID, "GetRecvFriendApplicationList callback: ", utils.structToJsonString(result))
	}()
}

func (u *UserRelated) GetSendFriendApplicationList(callback Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetSendFriendApplicationList args: ")
		result := u.getSendFriendApplicationList(callback, operationID)
		callback.OnSuccess(utils.structToJsonString(result))
		log.NewInfo(operationID, "GetSendFriendApplicationList callback: ", utils.structToJsonString(result))
	}()
}

func (u *UserRelated) AcceptFriendApplication(callback Base, params string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "AcceptFriendApplication args: ", params)
		var unmarshalParams sdk_params_callback.ProcessFriendApplicationParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback, operationID)
		u.processFriendApplication(callback, unmarshalParams, 1, operationID)
		callback.OnSuccess(utils.structToJsonString(sdk_params_callback.ProcessFriendApplicationCallback))
		log.NewInfo(operationID, "AcceptFriendApplication callback: ", utils.structToJsonString(sdk_params_callback.ProcessFriendApplicationCallback))
	}()
}

func (u *UserRelated) RefuseFriendApplication(callback Base, params string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "RefuseFriendApplication args: ", params)
		var unmarshalParams sdk_params_callback.ProcessFriendApplicationParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback, operationID)
		u.processFriendApplication(callback, unmarshalParams, -1, operationID)
		callback.OnSuccess(utils.structToJsonString(sdk_params_callback.ProcessFriendApplicationCallback))
		log.NewInfo(operationID, "RefuseFriendApplication callback: ", utils.structToJsonString(sdk_params_callback.ProcessFriendApplicationCallback))
	}()
}

func (u *UserRelated) CheckFriend(callback Base, params string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "CheckFriend args: ", params)
		var unmarshalParams sdk_params_callback.CheckFriendParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback, operationID)
		result := u.checkFriend(callback, unmarshalParams, operationID)
		callback.OnSuccess(utils.structToJsonString(result))
		log.NewInfo(operationID, "CheckFriend callback: ", utils.structToJsonString(result))
	}()
}

func (u *UserRelated) DeleteFriend(callback Base, friendUserID string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), friendUserID)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "DeleteFriend args: ", friendUserID)
		u.deleteFriend(friendUserID, callback, operationID)
		callback.OnSuccess(utils.structToJsonString(sdk_params_callback.DeleteFriendCallback))
		log.NewInfo(operationID, "DeleteFriend callback: ", utils.structToJsonString(sdk_params_callback.DeleteFriendCallback))
	}()
}

func (u *UserRelated) GetFriendList(callback Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetFriendList args: ")
		var filterLocalFriendList []*LocalFriend
		localFriendList, err := u._getAllFriendList()
		utils.checkErr(callback, err, operationID)
		localBlackUidList, err := u._getBlackListUid()
		utils.checkErr(callback, err, operationID)
		for _, v := range localFriendList {
			if !utils.isContain(v.FriendUserID, localBlackUidList) {
				filterLocalFriendList = append(filterLocalFriendList, v)
			}
		}
		callback.OnSuccess(utils.structToJsonString(filterLocalFriendList))
		log.NewInfo(operationID, "GetFriendList callback: ", utils.structToJsonString(filterLocalFriendList))
	}()
}

func (u *UserRelated) SetFriendRemark(callback Base, params string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "SetFriendRemark args: ", params)
		var unmarshalParams sdk_params_callback.SetFriendRemarkParams
		u.jsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback, operationID)
		u.setFriendRemark(unmarshalParams, callback, operationID)
		callback.OnSuccess(utils.structToJsonString(sdk_params_callback.SetFriendRemarkCallback))
		log.NewInfo(operationID, "SetFriendRemark callback: ", utils.structToJsonString(sdk_params_callback.SetFriendRemarkCallback))
	}()
}

func (u *UserRelated) AddBlack(callback Base, blackUserID, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), blackUserID)
	go func() {
		log.NewInfo(operationID, "AddToBlackList args: ", blackUserID)
		u.addBlack(callback, blackUserID, operationID)
		callback.OnSuccess("")
		log.NewInfo(operationID, "AddToBlackList callback: ")
	}()
}

func (u *UserRelated) GetBlackList(callback Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	go func() {
		log.NewInfo(operationID, "GetBlackList args: ")
		localBlackList, err := u._getBlackList()
		utils.checkErr(callback, err, operationID)
		callback.OnSuccess(utils.structToJsonString(localBlackList))
		log.NewInfo(operationID, "GetBlackList callback: ", utils.structToJsonString(localBlackList))
	}()
}

func (u *UserRelated) RemoveBlack(callback Base, removeUserID, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), removeUserID)
	go func() {
		log.NewInfo(operationID, "RemoveBlack args: ", removeUserID)
		u.removeBlack(callback, removeUserID, operationID)
		callback.OnSuccess("")
		log.NewInfo(operationID, "RemoveBlack callback")
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
