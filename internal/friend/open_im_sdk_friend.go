package friend

import (
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
)

func (f *Friend) GetDesignatedFriendsInfo(callback common.Base, params string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetDesignatedFriendsInfo args: ", params)
		var unmarshalList sdk_params_callback.GetDesignatedFriendsInfoParams
		common.JsonUnmarshal(params, &unmarshalList, callback, operationID)
		result := f.getDesignatedFriendsInfo(callback, unmarshalList, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "GetDesignatedFriendsInfo callback: ", utils.StructToJsonString(result), result)
	}()
}

func (f *Friend) AddFriend(callback common.Base, params string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "AddFriend args: ", params)
		var unmarshalAddFriendParams sdk_params_callback.AddFriendParams
		common.JsonUnmarshalAndArgsValidate(params, &unmarshalAddFriendParams, callback, operationID)
		f.addFriend(callback, unmarshalAddFriendParams, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.AddFriendCallback))
		log.NewInfo(operationID, "AddFriend callback: ", utils.StructToJsonString(sdk_params_callback.AddFriendCallback))
	}()
}

func (f *Friend) GetRecvFriendApplicationList(callback common.Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetRecvFriendApplicationList args: ")
		result := f.getRecvFriendApplicationList(callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(result))
		log.NewInfo(operationID, "GetRecvFriendApplicationList callback: ", utils.StructToJsonString(result))
	}()
}

func (f *Friend) GetSendFriendApplicationList(callback common.Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetSendFriendApplicationList args: ")
		result := f.getSendFriendApplicationList(callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(result))
		log.NewInfo(operationID, "GetSendFriendApplicationList callback: ", utils.StructToJsonString(result))
	}()
}

func (f *Friend) AcceptFriendApplication(callback common.Base, params string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "AcceptFriendApplication args: ", params)
		var unmarshalParams sdk_params_callback.ProcessFriendApplicationParams
		common.JsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback, operationID)
		f.processFriendApplication(callback, unmarshalParams, 1, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.ProcessFriendApplicationCallback))
		log.NewInfo(operationID, "AcceptFriendApplication callback: ", utils.StructToJsonString(sdk_params_callback.ProcessFriendApplicationCallback))
	}()
}

func (f *Friend) RefuseFriendApplication(callback common.Base, params string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "RefuseFriendApplication args: ", params)
		var unmarshalParams sdk_params_callback.ProcessFriendApplicationParams
		common.JsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback, operationID)
		f.processFriendApplication(callback, unmarshalParams, -1, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.ProcessFriendApplicationCallback))
		log.NewInfo(operationID, "RefuseFriendApplication callback: ", utils.StructToJsonString(sdk_params_callback.ProcessFriendApplicationCallback))
	}()
}

func (f *Friend) CheckFriend(callback common.Base, params string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "CheckFriend args: ", params)
		var unmarshalParams sdk_params_callback.CheckFriendParams
		common.JsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback, operationID)
		result := f.checkFriend(callback, unmarshalParams, operationID)
		callback.OnSuccess(utils.StructToJsonString(result))
		log.NewInfo(operationID, "CheckFriend callback: ", utils.StructToJsonString(result))
	}()
}

func (f *Friend) DeleteFriend(callback common.Base, friendUserID string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), friendUserID)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "DeleteFriend args: ", friendUserID)
		f.deleteFriend(friendUserID, callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.DeleteFriendCallback))
		log.NewInfo(operationID, "DeleteFriend callback: ", utils.StructToJsonString(sdk_params_callback.DeleteFriendCallback))
	}()
}

func (f *Friend) GetFriendList(callback common.Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetFriendList args: ")
		var filterLocalFriendList []*db.LocalFriend
		localFriendList, err := f.db.GetAllFriendList()
		common.CheckErr(callback, err, operationID)
		localBlackUidList, err := f.db.GetBlackListUid()
		common.CheckErr(callback, err, operationID)
		for _, v := range localFriendList {
			if !utils.IsContain(v.FriendUserID, localBlackUidList) {
				filterLocalFriendList = append(filterLocalFriendList, v)
			}
		}
		callback.OnSuccess(utils.StructToJsonString(filterLocalFriendList))
		log.NewInfo(operationID, "GetFriendList callback: ", utils.StructToJsonString(filterLocalFriendList))
	}()
}

func (f *Friend) SetFriendRemark(callback common.Base, params string, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), params)
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "SetFriendRemark args: ", params)
		var unmarshalParams sdk_params_callback.SetFriendRemarkParams
		common.JsonUnmarshalAndArgsValidate(params, &unmarshalParams, callback, operationID)
		f.setFriendRemark(unmarshalParams, callback, operationID)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.SetFriendRemarkCallback))
		log.NewInfo(operationID, "SetFriendRemark callback: ", utils.StructToJsonString(sdk_params_callback.SetFriendRemarkCallback))
	}()
}

func (f *Friend) AddBlack(callback common.Base, blackUserID, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), blackUserID)
	go func() {
		log.NewInfo(operationID, "AddToBlackList args: ", blackUserID)
		f.addBlack(callback, blackUserID, operationID)
		callback.OnSuccess("")
		log.NewInfo(operationID, "AddToBlackList callback: ")
	}()
}

func (f *Friend) GetBlackList(callback common.Base, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName())
	go func() {
		log.NewInfo(operationID, "GetBlackList args: ")
		localBlackList, err := f.db.GetBlackList()
		common.CheckErr(callback, err, operationID)
		callback.OnSuccess(utils.StructToJsonString(localBlackList))
		log.NewInfo(operationID, "GetBlackList callback: ", utils.StructToJsonString(localBlackList))
	}()
}

func (f *Friend) RemoveBlack(callback common.Base, removeUserID, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), removeUserID)
	go func() {
		log.NewInfo(operationID, "RemoveBlack args: ", removeUserID)
		f.removeBlack(callback, removeUserID, operationID)
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

func (f *Friend) SetFriendListener(listener OnFriendshipListener) bool {
	if listener == nil {
		return false
	}
	f.friendListener = listener
	return true
}

//func (f *Friend) ForceSyncFriendApplication() {
//	f.syncFriendApplication()
//}

//func (f *Friend) ForceSyncSelfFriendApplication() {
//	f.syncSelfFriendApplication()
//}

//func (f *Friend) ForceSyncFriend() {
//	f.syncFriendList()
//}

//func (f *Friend) ForceSyncBlackList() {
//	f.syncBlackList()
//}
