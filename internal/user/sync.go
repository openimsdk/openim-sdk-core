package user

import (
	"context"
	"fmt"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"

	"github.com/google/go-cmp/cmp"
)

func (u *User) SyncLoginUserInfo(ctx context.Context) error {
	remoteUser, err := u.GetSingleUserFromSvr(ctx, u.loginUserID)
	if err != nil {
		return err
	}
	localUser, err := u.GetLoginUser(ctx, u.loginUserID)
	if err != nil {
		return err
	}
	return u.userSyncer.Sync(ctx, []*model_struct.LocalUser{remoteUser}, []*model_struct.LocalUser{localUser}, nil)
}

// deprecated
func (u *User) SyncLoginUserInfoOld(ctx context.Context) error {
	// log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svr, err := u.GetSelfUserInfoFromSvr(ctx)
	if err != nil {
		// log.Error(operationID, "GetSelfUserInfoFromSvr failed", err.Error())
		return err
	}

	onServer := common.TransferToLocalUserInfo(svr)
	onLocal, err := u.GetLoginUser(u.loginUserID)
	if err != nil {
		// log.Warn(operationID, "GetLoginUser failed ", err.Error())
		onLocal = &model_struct.LocalUser{}
	}
	if !cmp.Equal(onServer, onLocal) {
		if onLocal.UserID == "" {
			if err = u.InsertLoginUser(onServer); err != nil {
				log.Error(operationID, "InsertLoginUser failed ", *onServer, err.Error())
				return
			}

		} else {
			err = u.UpdateLoginUserByMap(onServer, map[string]interface{}{"name": onServer.Nickname, "face_url": onServer.FaceURL,
				"gender": onServer.Gender, "phone_number": onServer.PhoneNumber, "birth_time": onServer.BirthTime, "email": onServer.Email, "create_time": onServer.CreateTime, "app_manger_level": onServer.AppMangerLevel, "ex": onServer.Ex, "attached_info": onServer.AttachedInfo, "global_recv_msg_opt": onServer.GlobalRecvMsgOpt})
			fmt.Println("UpdateLoginUser ", *onServer, svr)
			if err != nil {
				log.Error(operationID, "UpdateLoginUser failed ", *onServer, err.Error())
				return
			}
		}
		callbackData := sdk.SelfInfoUpdatedCallback(*onServer)
		if u.listener == nil {
			log.Error(operationID, "u.listener == nil")
			return
		}
		u.listener.OnSelfInfoUpdated(utils.StructToJsonString(callbackData))
		log.Info(operationID, "OnSelfInfoUpdated", utils.StructToJsonString(callbackData))
		if onLocal.Nickname == onServer.Nickname && onLocal.FaceURL == onServer.FaceURL {
			log.NewInfo(operationID, "OnSelfInfoUpdated nickname faceURL unchanged", callbackData)
			return
		}
		_ = common.TriggerCmdUpdateMessage(common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName, Args: common.UpdateMessageInfo{UserID: callbackData.UserID, FaceURL: callbackData.FaceURL, Nickname: callbackData.Nickname}}, u.conversationCh)

	}
	return nil
}
