package user

import (
	"github.com/mitchellh/mapstructure"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

type User struct {
	*db.DataBase
	p           *ws.PostApi
	loginUserID string
	listener    OnUserListener
}

func (u *User) SetListener(listener OnUserListener) {
	u.listener = listener
}

func NewUser(dataBase *db.DataBase, p *ws.PostApi, loginUserID string) *User {
	return &User{DataBase: dataBase, p: p, loginUserID: loginUserID}
}

type OnUserListener interface {
	OnSelfInfoUpdated(userInfo string)
}

func (u *User) DoNotification(msg *api.MsgData) {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg)
	if u.listener == nil {
		log.Error(operationID, "listener == nil")
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
	u.SyncLoginUserInfo(operationID)
}

func (u *User) SyncLoginUserInfo(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svr, err := u.GetSelfUserInfoFromSvr(operationID)
	if err != nil {
		log.Error(operationID, "GetSelfUserInfoFromSvr failed")
		return
	}
	onServer := common.TransferToLocalUserInfo(svr)
	onLocal, err := u.GetLoginUser()
	if err != nil {
		log.Error(operationID, "TransferToLocalUserInfo failed")
		return
	}
	if onServer != onLocal {
		u.UpdateLoginUser(onServer)
		if err != nil {
			log.Error(operationID, "UpdateLoginUser failed", onServer)
			return
		}
		callbackData := sdk.SelfInfoUpdatedCallback(*onServer)
		if u.listener == nil {
			log.Error(operationID, "u.listener == nil")
			return
		}
		u.listener.OnSelfInfoUpdated(utils.StructToJsonString(callbackData))
	}
}

func (u *User) GetUsersInfoFromSvr(callback common.Base, UserIDList sdk.GetUsersInfoParam, operationID string) sdk.GetUsersInfoCallback {
	apiReq := api.GetUsersInfoReq{}
	apiReq.OperationID = operationID
	apiReq.UserIDList = UserIDList
	commData := u.p.PostFatalCallback(callback, constant.GetUsersInfoRouter, apiReq, apiReq.OperationID)
	apiResp := make([]*api.PublicUserInfo, 0)
	common.MapstructureDecode(commData.Data, &apiResp, callback, apiReq.OperationID)
	return apiResp
}

func (u *User) getSelfUserInfo(callback common.Base, operationID string) sdk.GetSelfUserInfoCallback {
	userInfo, err := u.GetLoginUser()
	if err != nil {
		callback.OnError(constant.ErrDB.ErrCode, constant.ErrDB.ErrMsg)
	}
	return userInfo
}

func (u *User) updateSelfUserInfo(callback common.Base, userInfo sdk.SetSelfUserInfoParam, operationID string) *api.CommDataResp {
	apiReq := api.UpdateSelfUserInfoReq{}
	apiReq.OperationID = operationID
	apiReq.UserInfo = api.UserInfo(userInfo)
	apiReq.UserID = u.loginUserID
	commData := u.p.PostFatalCallback(callback, constant.UpdateSelfUserInfoRouter, apiReq, apiReq.OperationID)
	apiResp := api.CommDataResp{}
	common.MapstructureDecode(commData.Data, &apiResp, callback, apiReq.OperationID)
	return &apiResp
}

func (u *User) GetSelfUserInfoFromSvr(operationID string) (*api.UserInfo, error) {
	log.Debug(operationID, utils.GetSelfFuncName())
	apiReq := api.GetSelfUserInfoReq{}
	apiReq.OperationID = operationID
	apiReq.UserID = u.loginUserID
	commData, err := u.p.PostReturn(constant.GetSelfUserInfoRouter, apiReq)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	apiResp := api.UserInfo{}
	mapstructure.Decode(commData.Data, &apiResp)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return &apiResp, nil
}

func (u *User) DoUserNotification(msg *api.MsgData) {
	if msg.SendID == u.loginUserID && msg.SenderPlatformID == sdk_struct.SvrConf.Platform {
		return
	}
}
