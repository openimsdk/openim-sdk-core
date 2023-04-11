package friend

import (
	"context"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	sdk "open_im_sdk/pkg/sdk_params_callback"
)

func (f *Friend) SyncSelfFriendApplication(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := f.getSelfFriendApplicationFromServer(operationID)
	if err != nil {
		log.NewError(operationID, "getSelfFriendApplicationFromServer failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalFriendRequest(svrList)
	onLocal, err := f.db.GetSendFriendApplication()
	if err != nil {
		log.NewError(operationID, "GetSendFriendApplication failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "list", svrList, onServer, onLocal)

	aInBNot, bInANot, sameA, sameB := common.CheckFriendRequestDiff(onServer, onLocal)
	log.Debug(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := f.db.InsertFriendRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertFriendRequest failed ", err.Error())
			continue
		}
		callbackData := sdk.FriendApplicationAddedCallback(*onServer[index])
		if f.friendListener != nil {
			f.friendListener.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnFriendApplicationAdded", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range sameA {
		err := f.db.UpdateFriendRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateFriendRequest failed ", err.Error(), *onServer[index])
			continue
		} else {
			if onServer[index].HandleResult == constant.FriendResponseRefuse {
				callbackData := sdk.FriendApplicationRejectCallback(*onServer[index])
				if f.friendListener != nil {

					f.friendListener.OnFriendApplicationRejected(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnFriendApplicationRejected", utils.StructToJsonString(callbackData))
				}

			} else if onServer[index].HandleResult == constant.FriendResponseAgree {
				callbackData := sdk.FriendApplicationAcceptCallback(*onServer[index])
				if f.friendListener != nil {
					f.friendListener.OnFriendApplicationAccepted(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnFriendApplicationAccepted", utils.StructToJsonString(callbackData))
				}
				if f.listenerForService != nil {
					f.listenerForService.OnFriendApplicationAccepted(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnFriendApplicationAccepted", utils.StructToJsonString(callbackData))
				}
			} else {
				callbackData := sdk.FriendApplicationAddedCallback(*onServer[index])
				if f.friendListener != nil {
					f.friendListener.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnFriendApplicationAdded", utils.StructToJsonString(callbackData))
				}
			}
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriendRequestBothUserID(onLocal[index].FromUserID, onLocal[index].ToUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriendRequestBothUserID failed ", err.Error(), onLocal[index].FromUserID, onLocal[index].ToUserID)
			continue
		}
		callbackData := sdk.FriendApplicationDeletedCallback(*onLocal[index])
		if f.friendListener != nil {

			f.friendListener.OnFriendApplicationDeleted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnFriendApplicationDeleted", utils.StructToJsonString(callbackData))
		}
	}
}

// recv
func (f *Friend) SyncFriendApplication(ctx context.Context) error {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := f.getFriendApplicationFromServer(operationID)
	if err != nil {
		log.NewError(operationID, "getFriendApplicationFromServer failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalFriendRequest(svrList)
	onLocal, err := f.db.GetRecvFriendApplication()
	if err != nil {
		log.NewError(operationID, "GetRecvFriendApplication failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "list", svrList, onServer, onLocal)

	aInBNot, bInANot, sameA, sameB := common.CheckFriendRequestDiff(onServer, onLocal)
	log.Debug(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := f.db.InsertFriendRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertFriendRequest failed ", err.Error())
			continue
		}
		callbackData := sdk.FriendApplicationAddedCallback(*onServer[index])
		//f.friendListener.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
		if f.friendListener != nil {
			f.friendListener.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnReceiveFriendApplicationAdded", utils.StructToJsonString(callbackData))
		}
		if f.listenerForService != nil {
			f.listenerForService.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnReceiveFriendApplicationAdded", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range sameA {
		err := f.db.UpdateFriendRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateFriendRequest failed ", err.Error(), *onServer[index])
			continue
		} else {
			if onServer[index].HandleResult == constant.FriendResponseRefuse {
				callbackData := sdk.FriendApplicationRejectCallback(*onServer[index])
				if f.friendListener != nil {

					f.friendListener.OnFriendApplicationRejected(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnFriendApplicationRejected", utils.StructToJsonString(callbackData))
				}
			} else if onServer[index].HandleResult == constant.FriendResponseAgree {
				callbackData := sdk.FriendApplicationAcceptCallback(*onServer[index])
				if f.friendListener != nil {

					f.friendListener.OnFriendApplicationAccepted(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnFriendApplicationAccepted", utils.StructToJsonString(callbackData))
				}
			} else {
				callbackData := sdk.FriendApplicationAddedCallback(*onServer[index])
				if f.friendListener != nil {
					f.friendListener.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnReceiveFriendApplicationAdded", utils.StructToJsonString(callbackData))
				}
				if f.listenerForService != nil {
					f.listenerForService.OnFriendApplicationAdded(utils.StructToJsonString(callbackData))
					log.Info(operationID, "OnReceiveFriendApplicationAdded", utils.StructToJsonString(callbackData))
				}
			}
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriendRequestBothUserID(onLocal[index].FromUserID, onLocal[index].ToUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriendRequestBothUserID failed ", err.Error(), onLocal[index].FromUserID, onLocal[index].ToUserID)
			continue
		}
		callbackData := sdk.FriendApplicationDeletedCallback(*onLocal[index])
		if f.friendListener != nil {

			f.friendListener.OnFriendApplicationDeleted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnReceiveFriendApplicationDeleted", utils.StructToJsonString(callbackData))
		}
	}
}

func (f *Friend) SyncFriendList(ctx context.Context) error {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := f.getServerFriendList(operationID)
	if err != nil {
		log.NewError(operationID, "getServerFriendList failed ", err.Error())
		return
	}
	friendsInfoOnServer := common.TransferToLocalFriend(svrList)
	friendsInfoOnLocal, err := f.db.GetAllFriendList()
	if err != nil {
		log.NewError(operationID, "_getAllFriendList failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "list ", svrList, friendsInfoOnServer, friendsInfoOnLocal)
	for _, v := range friendsInfoOnServer {
		log.NewDebug(operationID, "friendsInfoOnServer ", *v)
	}
	aInBNot, bInANot, sameA, sameB := common.CheckFriendListDiff(friendsInfoOnServer, friendsInfoOnLocal)
	log.NewInfo(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := f.db.InsertFriend(friendsInfoOnServer[index])
		if err != nil {
			log.NewError(operationID, "_insertFriend failed ", err.Error())
			continue
		}
		callbackData := sdk.FriendAddedCallback(*friendsInfoOnServer[index])
		if f.friendListener != nil {

			f.friendListener.OnFriendAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnFriendAdded", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range sameA {
		callbackData := sdk.FriendInfoChangedCallback(*friendsInfoOnServer[index])
		localFriend, err := f.db.GetFriendInfoByFriendUserID(callbackData.FriendUserID)
		if err != nil {
			log.NewError(operationID, "GetFriendInfoByFriendUserID failed ", err.Error(), "userID", callbackData.FriendUserID)
			continue
		}
		err = f.db.UpdateFriend(friendsInfoOnServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateFriendRequest failed ", err.Error(), *friendsInfoOnServer[index])
			continue
		} else {
			callbackData := sdk.FriendInfoChangedCallback(*friendsInfoOnServer[index])
			if f.friendListener != nil {
				f.friendListener.OnFriendInfoChanged(utils.StructToJsonString(callbackData))
				if localFriend.Nickname == callbackData.Nickname && localFriend.FaceURL == callbackData.FaceURL && localFriend.Remark == callbackData.Remark {
					log.NewInfo(operationID, "OnFriendInfoChanged nickname faceURL unchanged", callbackData.FriendUserID, localFriend.Nickname, localFriend.FaceURL)
					continue
				}
				if callbackData.Remark != "" {
					callbackData.Nickname = callbackData.Remark
				}
				common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName, Args: common.SourceIDAndSessionType{SourceID: callbackData.FriendUserID, SessionType: constant.SingleChatType}}, f.conversationCh)
				common.TriggerCmdUpdateMessage(common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName, Args: common.UpdateMessageInfo{UserID: callbackData.FriendUserID, FaceURL: callbackData.FaceURL, Nickname: callbackData.Nickname}}, f.conversationCh)
				log.Info(operationID, "OnFriendInfoChanged", utils.StructToJsonString(callbackData))
			}
		}
	}
	for _, index := range bInANot {
		err := f.db.DeleteFriendDB(friendsInfoOnLocal[index].FriendUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriend failed ", err.Error())
			continue
		}
		callbackData := sdk.FriendDeletedCallback(*friendsInfoOnLocal[index])
		if f.friendListener != nil {

			f.friendListener.OnFriendDeleted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnFriendDeleted", utils.StructToJsonString(callbackData))
		}
	}
}

func (f *Friend) SyncBlackList(ctx context.Context) error {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := f.getServerBlackList(operationID)
	if err != nil {
		log.NewError(operationID, "getServerBlackList failed ", err.Error())
		return
	}
	blackListOnServer := common.TransferToLocalBlack(svrList, f.loginUserID)
	blackListOnLocal, err := f.db.GetBlackListDB()
	if err != nil {
		log.NewError(operationID, "_getBlackList failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "list ", svrList, blackListOnServer, blackListOnLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckBlackListDiff(blackListOnServer, blackListOnLocal)
	log.NewInfo(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := f.db.InsertBlack(blackListOnServer[index])
		if err != nil {
			log.NewError(operationID, "_insertFriend failed ", err.Error())
			continue
		}
		callbackData := sdk.BlackAddCallback(*blackListOnServer[index])
		if f.friendListener != nil {

			f.friendListener.OnBlackAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnBlackAdded", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range sameA {
		err := f.db.UpdateBlack(blackListOnServer[index])
		if err != nil {
			log.NewError(operationID, "_updateFriend failed ", err.Error())
			continue
		}
		//todo : add black info update callback
		log.Info(operationID, "black info update, do nothing ", blackListOnServer[index])
	}
	for _, index := range bInANot {
		err := f.db.DeleteBlack(blackListOnLocal[index].BlockUserID)
		if err != nil {
			log.NewError(operationID, "_deleteFriend failed ", err.Error())
			continue
		}
		callbackData := sdk.BlackDeletedCallback(*blackListOnLocal[index])
		if f.friendListener != nil {
			f.friendListener.OnBlackDeleted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnBlackDeleted", utils.StructToJsonString(callbackData))
		}
	}
}
