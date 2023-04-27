package user

import (
	"fmt"
	comm "open_im_sdk/internal/common"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"

	"github.com/google/go-cmp/cmp"

	//"github.com/mitchellh/mapstructure"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	sdkstruct "open_im_sdk/sdk_struct"
)

// User is a struct that represents a user in the system.
type User struct {
	db_interface.DataBase
	p              *ws.PostApi
	loginUserID    string
	listener       open_im_sdk_callback.OnUserListener
	loginTime      int64
	conversationCh chan common.Cmd2Value
}

// LoginTime gets the login time of the user.
func (u *User) LoginTime() int64 {
	return u.loginTime
}

// SetLoginTime sets the login time of the user.
func (u *User) SetLoginTime(loginTime int64) {
	u.loginTime = loginTime
}

// SetListener sets the user's listener.
func (u *User) SetListener(listener open_im_sdk_callback.OnUserListener) {
	u.listener = listener
}

// NewUser creates a new User object.
func NewUser(dataBase db_interface.DataBase, p *ws.PostApi, loginUserID string, conversationCh chan common.Cmd2Value) *User {
	return &User{DataBase: dataBase, p: p, loginUserID: loginUserID, conversationCh: conversationCh}
}

// DoNotification handles incoming notifications for the user.
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
			if err := u.userInfoUpdatedNotification(msg, operationID); err != nil {
				log.Error(operationID, "userInfoUpdatedNotification failed ", err.Error())
			}
		default:
			log.Error(operationID, "type failed ", msg.ClientMsgID, msg.ServerMsgID, msg.ContentType)
		}
	}()
}

// userInfoUpdatedNotification handles notifications about updated user information.
func (u *User) userInfoUpdatedNotification(msg *api.MsgData, operationID string) error {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.UserInfoUpdatedTips{}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "comm.UnmarshalTips failed ", err.Error(), msg.Content)
		return err
	}
	if detail.UserID == u.loginUserID {
		log.Info(operationID, "detail.UserID == u.loginUserID, SyncLoginUserInfo", detail.UserID)
		u.SyncLoginUserInfo(operationID)
	}
	log.Info(operationID, "detail.UserID != u.loginUserID, do nothing", detail.UserID, u.loginUserID)
	return nil
}

// SyncLoginUserInfo synchronizes the user's information from the server.
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

// GetUsersInfoFromSvr retrieves user information from the server.
func (u *User) GetUsersInfoFromSvr(callback open_im_sdk_callback.Base, UserIDList sdk.GetUsersInfoParam, operationID string) []*api.PublicUserInfo {
	apiReq := api.GetUsersInfoReq{OperationID: operationID, UserIDList: UserIDList}
	apiResp := api.GetUsersInfoResp{}
	u.p.PostFatalCallback(callback, constant.GetUsersInfoRouter, apiReq, &apiResp.UserInfoList, apiReq.OperationID)
	return apiResp.UserInfoList
}

// GetUsersInfoFromSvrNoCallback retrieves user information from the server.
func (u *User) GetUsersInfoFromSvrNoCallback(UserIDList sdk.GetUsersInfoParam, operationID string) ([]*api.PublicUserInfo, error) {
	apiReq := api.GetUsersInfoReq{OperationID: operationID, UserIDList: UserIDList}
	apiResp := api.GetUsersInfoResp{}
	err := u.p.PostReturn(constant.GetUsersInfoRouter, apiReq, &apiResp.UserInfoList)
	return apiResp.UserInfoList, err
}

// GetUsersInfoFromCacheSvr retrieves user information from the cache server.
func (u *User) GetUsersInfoFromCacheSvr(UserIDList sdk.GetUsersInfoParam, operationID string) ([]*api.PublicUserInfo, error) {
	apiReq := api.GetUsersInfoReq{OperationID: operationID, UserIDList: UserIDList}
	apiResp := api.GetUsersInfoResp{}
	err := u.p.PostReturn(constant.GetUsersInfoFromCacheRouter, apiReq, &apiResp.UserInfoList)
	return apiResp.UserInfoList, err
}

// getSelfUserInfo retrieves the user's information.
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

// updateSelfUserInfo updates the user's information.
func (u *User) updateSelfUserInfo(callback open_im_sdk_callback.Base, userInfo sdk.SetSelfUserInfoParam, operationID string) {
	apiReq := api.UpdateSelfUserInfoReq{}
	apiReq.OperationID = operationID
	apiReq.ApiUserInfo = api.ApiUserInfo(userInfo)
	apiReq.UserID = u.loginUserID
	u.p.PostFatalCallback(callback, constant.UpdateSelfUserInfoRouter, apiReq, nil, apiReq.OperationID)
	u.SyncLoginUserInfo(operationID)
}

// GetSelfUserInfoFromSvr retrieves the user's information from the server.
func (u *User) GetSelfUserInfoFromSvr(operationID string) (*api.UserInfo, error) {
	apiReq := api.GetSelfUserInfoReq{OperationID: operationID, UserID: u.loginUserID}
	apiResp := api.GetSelfUserInfoResp{UserInfo: &api.UserInfo{}}
	err := u.p.PostReturn(constant.GetSelfUserInfoRouter, apiReq, &apiResp.UserInfo)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return apiResp.UserInfo, nil
}

// DoUserNotification handles incoming notifications for the user.
func (u *User) DoUserNotification(msg *api.MsgData) {
	if msg.SendID == u.loginUserID && msg.SenderPlatformID == sdkstruct.SvrConf.Platform {
		return
	}
}

// ParseTokenFromSvr parses a token from the server.
func (u *User) ParseTokenFromSvr(operationID string) (uint32, error) {
	apiReq := api.ParseTokenReq{OperationID: operationID}
	apiResp := api.ParseTokenResp{}
	err := u.p.PostReturn(constant.ParseTokenRouter, apiReq, &apiResp.ExpireTime)
	if err != nil {
		return 0, utils.Wrap(err, apiReq.OperationID)
	}
	log.Info(operationID, "apiResp.ExpireTime.ExpireTimeSeconds ", apiResp.ExpireTime)
	return apiResp.ExpireTime.ExpireTimeSeconds, nil
}
