package user

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	comm "open_im_sdk/internal/common"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"

	//"github.com/mitchellh/mapstructure"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

type User struct {
	db_interface.DataBase
	p              *ws.PostApi
	loginUserID    string
	listener       open_im_sdk_callback.OnUserListener
	loginTime      int64
	conversationCh chan common.Cmd2Value
}

func (u *User) LoginTime() int64 {
	return u.loginTime
}

func (u *User) SetLoginTime(loginTime int64) {
	u.loginTime = loginTime
}

func (u *User) SetListener(listener open_im_sdk_callback.OnUserListener) {
	u.listener = listener
}

func NewUser(dataBase db_interface.DataBase, p *ws.PostApi, loginUserID string, conversationCh chan common.Cmd2Value) *User {
	return &User{DataBase: dataBase, p: p, loginUserID: loginUserID, conversationCh: conversationCh}
}

func (u *User) DoNotification(msg *api.MsgData) {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg)
	if u.listener == nil {
		log.Error(operationID, "listener == nil")
		return
	}

	if msg.SendTime < u.loginTime {
		log.Warn(operationID, "ignore notification ", msg.ClientMsgID, msg.ServerMsgID, msg.Seq, msg.ContentType)
		return
	}
	go func() {
		switch msg.ContentType {
		case constant.UserInfoUpdatedNotification:
			u.userInfoUpdatedNotification(msg, operationID)
		default:
			log.Error(operationID, "type failed ", msg.ClientMsgID, msg.ServerMsgID, msg.ContentType)
		}
	}()
}

func (u *User) userInfoUpdatedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	var detail api.UserInfoUpdatedTips
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg.Content)
		return
	}
	if detail.UserID == u.loginUserID {
		log.Info(operationID, "detail.UserID == u.loginUserID, SyncLoginUserInfo", detail.UserID)
		u.SyncLoginUserInfo(operationID)
	} else {
		log.Info(operationID, "detail.UserID != u.loginUserID, do nothing", detail.UserID, u.loginUserID)
	}
}

func (u *User) SyncLoginUserInfo(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svr, err := u.GetSelfUserInfoFromSvr(operationID)
	if err != nil {
		log.Error(operationID, "GetSelfUserInfoFromSvr failed", err.Error())
		return
	}
	onServer := common.TransferToLocalUserInfo(svr)
	onLocal, err := u.GetLoginUser(u.loginUserID)
	if err != nil {
		log.Warn(operationID, "GetLoginUser failed ", err.Error())
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
}

func (u *User) GetUsersInfoFromSvr(callback open_im_sdk_callback.Base, UserIDList sdk.GetUsersInfoParam, operationID string) []*api.PublicUserInfo {
	apiReq := api.GetUsersInfoReq{}
	apiReq.OperationID = operationID
	apiReq.UserIDList = UserIDList
	apiResp := api.GetUsersInfoResp{}
	u.p.PostFatalCallback(callback, constant.GetUsersInfoRouter, apiReq, &apiResp.UserInfoList, apiReq.OperationID)
	return apiResp.UserInfoList
}

func (u *User) GetUsersInfoFromSvrNoCallback(UserIDList sdk.GetUsersInfoParam, operationID string) ([]*api.PublicUserInfo, error) {
	apiReq := api.GetUsersInfoReq{}
	apiReq.OperationID = operationID
	apiReq.UserIDList = UserIDList
	apiResp := api.GetUsersInfoResp{}
	err := u.p.PostReturn(constant.GetUsersInfoRouter, apiReq, &apiResp.UserInfoList)
	return apiResp.UserInfoList, err
}

func (u *User) GetUsersInfoFromCacheSvr(UserIDList sdk.GetUsersInfoParam, operationID string) ([]*api.PublicUserInfo, error) {
	apiReq := api.GetUsersInfoReq{}
	apiReq.OperationID = operationID
	apiReq.UserIDList = UserIDList
	apiResp := api.GetUsersInfoResp{}
	err := u.p.PostReturn(constant.GetUsersInfoFromCacheRouter, apiReq, &apiResp.UserInfoList)
	return apiResp.UserInfoList, err
}

func (u *User) getSelfUserInfo(callback open_im_sdk_callback.Base, operationID string) sdk.GetSelfUserInfoCallback {
	userInfo, errLocal := u.GetLoginUser(u.loginUserID)
	if errLocal != nil {
		svr, errServer := u.GetSelfUserInfoFromSvr(operationID)
		if errServer != nil {
			log.Error(operationID, "GetSelfUserInfoFromSvr failed", errServer.Error())
			common.CheckDBErrCallback(callback, errServer, operationID)
		}
		userInfo = common.TransferToLocalUserInfo(svr)
	}
	return userInfo
}

func (u *User) updateSelfUserInfo(callback open_im_sdk_callback.Base, userInfo sdk.SetSelfUserInfoParam, operationID string) {
	apiReq := api.UpdateSelfUserInfoReq{}
	apiReq.OperationID = operationID
	apiReq.ApiUserInfo = api.ApiUserInfo(userInfo)
	apiReq.UserID = u.loginUserID
	u.p.PostFatalCallback(callback, constant.UpdateSelfUserInfoRouter, apiReq, nil, apiReq.OperationID)
	u.SyncLoginUserInfo(operationID)
}

func (u *User) GetSelfUserInfoFromSvr(operationID string) (*api.UserInfo, error) {
	log.Debug(operationID, utils.GetSelfFuncName())
	apiReq := api.GetSelfUserInfoReq{}
	apiReq.OperationID = operationID
	apiReq.UserID = u.loginUserID
	apiResp := api.GetSelfUserInfoResp{UserInfo: &api.UserInfo{}}
	err := u.p.PostReturn(constant.GetSelfUserInfoRouter, apiReq, &apiResp.UserInfo)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return apiResp.UserInfo, nil
}

func (u *User) DoUserNotification(msg *api.MsgData) {
	if msg.SendID == u.loginUserID && msg.SenderPlatformID == sdk_struct.SvrConf.Platform {
		return
	}
}

func (u *User) ParseTokenFromSvr(operationID string) (uint32, error) {
	apiReq := api.ParseTokenReq{}
	apiReq.OperationID = operationID
	apiResp := api.ParseTokenResp{}
	err := u.p.PostReturn(constant.ParseTokenRouter, apiReq, &apiResp.ExpireTime)
	if err != nil {
		return 0, utils.Wrap(err, apiReq.OperationID)
	}
	log.Info(operationID, "apiResp.ExpireTime.ExpireTimeSeconds ", apiResp.ExpireTime)
	return apiResp.ExpireTime.ExpireTimeSeconds, nil
}
