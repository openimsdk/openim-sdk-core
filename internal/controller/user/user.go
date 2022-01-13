package user

import (
	"github.com/mitchellh/mapstructure"
	ws "open_im_sdk/internal/controller/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

type User struct {
	*db.DataBase
	p           *ws.PostApi
	loginUserID string
}

func NewUser(dataBase *db.DataBase, p *ws.PostApi, loginUserID string) *User {
	return &User{DataBase: dataBase, p: p, loginUserID: loginUserID}
}

func (u *User) SyncLoginUserInfo() error {
	svr, err := u.GetSelfUserInfoFromSvr()
	if err != nil {
		return utils.Wrap(err, "_getSelfUserInfoFromSvr failed")
	}
	onServer := common.TransferToLocalUserInfo(svr)
	onLocal, err := u.GetLoginUser()
	if err != nil {
		return utils.Wrap(err, "GetLoginUser failed")
	}
	if onServer != onLocal {
		u.UpdateLoginUser(onServer)
	}
	return nil
}

func (u *User) getUsersInfoFromSvr(callback common.Base, UserIDList sdk.GetUsersInfoParam, operationID string) sdk.GetUsersInfoCallback {
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

func (u *User) GetSelfUserInfoFromSvr() (*api.UserInfo, error) {
	log.Debug("0", utils.GetSelfFuncName())

	apiReq := api.GetSelfUserInfoReq{}
	apiReq.OperationID = utils.OperationIDGenerator()
	apiReq.UserID = u.loginUserID
	commData, err := u.p.PostReturn(constant.GetSelfUserInfo, apiReq)
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
