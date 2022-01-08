package user

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"open_im_sdk/internal/open_im_sdk"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

func (u *User) SyncLoginUserInfo() error {
	userSvr, err := u.getServerUserInfo()
	if err != nil {
		log.NewError("0", "getServerUserInfo failed , user: ", err.Error())
		return err
	}

	log.NewInfo("0", "getServerUserInfo ", userSvr)

	userLocal, err := u._getLoginUser()
	needInsert := 0
	if err != nil {
		log.NewError("0", "_getLoginUser failed  ", err.Error())
		needInsert = 1
	}

	if utils.CompFields(&userLocal, &userSvr) {
		return nil
	}

	var updateLocalUser db.LocalUser
	copier.Copy(&updateLocalUser, userSvr)
	log.NewInfo("0", "copy: ", updateLocalUser)
	if needInsert == 1 {
		err = u._insertLoginUser(&updateLocalUser)
		if err != nil {
			log.NewError("0 ", "_insertLoginUser failed ", err.Error())
		}
		return err
	}
	err = u._updateLoginUser(&updateLocalUser)
	if err != nil {
		log.NewError("0 ", "_updateLoginUser failed ", err.Error())
	}
	return err

	//if err != nil {
	//	return err
	//}
	//sdkLog("getServerUserInfo ok, user: ", *userSvr)
	//
	//userLocal, err := u.getLoginUserInfoFromLocal()
	//userLocal, err := u._getLoginUser()
	//if err != nil {
	//	return err
	//}
	//sdkLog("getLoginUserInfoFromLocal ok, user: ", userLocal)
	//
	//if userSvr.Uid != userLocal.Uid ||
	//	userSvr.Name != userLocal.Name ||
	//	userSvr.Icon != userLocal.Icon ||
	//	userSvr.Gender != userLocal.Gender ||
	//	userSvr.Mobile != userLocal.Mobile ||
	//	userSvr.Birth != userLocal.Birth ||
	//	userSvr.Email != userLocal.Email ||
	//	userSvr.Ex != userLocal.Ex {
	//	bUserInfo, err := json.Marshal(userSvr)
	//	if err != nil {
	//		sdkLog("marshal failed, ", err.Error())
	//		return err
	//	}
	//
	//	copier.Copy(a, b)
	//	err = u._updateLoginUser(userSvr)
	//	if err != nil {
	//		u.cb.OnSelfInfoUpdated(string(bUserInfo))
	//	}
	//}
}

func (u *open_im_sdk.UserRelated) getServerUserInfo() (*server_api_params.UserInfo, error) {
	apiReq := server_api_params.GetUserInfoReq{OperationID: utils.operationIDGenerator(), UserIDList: []string{u.loginUserID}}
	resp, err := utils.post2Api(open_im_sdk.getUserInfoRouter, apiReq, u.token)
	commData, err := utils.checkErrAndRespReturn(err, resp, apiReq.OperationID)
	if err != nil {
		return nil, utils.wrap(err, apiReq.OperationID)
	}
	realData := server_api_params.GetUserInfoResp{}
	err = mapstructure.Decode(commData.Data, &realData.UserInfoList)
	if err != nil {
		log.NewError(apiReq.OperationID, "Decode failed ", err.Error())
		return nil, err
	}
	log.NewInfo(apiReq.OperationID, "realData.UserInfoList", realData.UserInfoList, commData.Data)
	if len(realData.UserInfoList) == 0 {
		log.NewInfo(apiReq.OperationID, "failed, no user : ", u.loginUserID)
		return nil, errors.New("no login user")
	}
	log.NewInfo(apiReq.OperationID, "realData.UserInfoList[0]", realData.UserInfoList[0])
	return realData.UserInfoList[0], nil
}

func (u *open_im_sdk.UserRelated) getUserNameAndFaceUrlByUid(uid string) (faceUrl, name string, err error) {
	friendInfo, err := u._getFriendInfoByFriendUserID(uid)
	if err != nil {
		return "", "", err
	}
	if friendInfo.FriendUserID == "" {
		userInfo, err := u.getUserInfoByUid(uid)
		if err != nil {
			return "", "", err
		} else {
			return userInfo.Icon, userInfo.Name, nil
		}
	} else {
		if friendInfo.Remark != "" {
			return friendInfo.FaceUrl, friendInfo.Remark, nil
		} else {
			return friendInfo.FaceUrl, friendInfo.Nickname, nil
		}
	}
}

func (u *open_im_sdk.UserRelated) getUserInfoByUid(uid string) (*open_im_sdk.userInfo, error) {
	var uidList []string
	uidList = append(uidList, uid)
	resp, err := utils.post2Api(open_im_sdk.getUserInfoRouter, open_im_sdk.paramsGetUserInfo{OperationID: utils.operationIDGenerator(), UidList: uidList}, u.token)
	if err != nil {
		utils.sdkLog("post2Api failed, ", open_im_sdk.getUserInfoRouter, uidList, err.Error())
		return nil, err
	}
	utils.sdkLog("post api: ", open_im_sdk.getUserInfoRouter, open_im_sdk.paramsGetUserInfo{OperationID: utils.operationIDGenerator(), UidList: uidList}, "uid ", uid)
	var userResp open_im_sdk.getUserInfoResp
	err = json.Unmarshal(resp, &userResp)
	if err != nil {
		utils.sdkLog("Unmarshal failed, ", resp, err.Error())
		return nil, err
	}

	if userResp.ErrCode != 0 {
		utils.sdkLog("errcode: ", userResp.ErrCode, "errmsg:", userResp.ErrMsg)
		return nil, errors.New(userResp.ErrMsg)
	}

	if len(userResp.Data) == 0 {
		utils.sdkLog("failed, no user :", uid)
		return nil, errors.New("no user")
	}
	return &userResp.Data[0], nil
}

func (u *open_im_sdk.UserRelated) setSelfInfo(msg *server_api_params.MsgData) {
	var uidList []string
	uidList = append(uidList, msg.SendID)
	resp, err := utils.post2Api(open_im_sdk.getUserInfoRouter, open_im_sdk.paramsGetUserInfo{OperationID: utils.operationIDGenerator(), UidList: uidList}, u.token)
	if err != nil {
		utils.sdkLog("post2Api failed, ", open_im_sdk.getUserInfoRouter, uidList, err.Error())
		return
	}
	var userResp open_im_sdk.getUserInfoResp
	err = json.Unmarshal(resp, &userResp)
	if err != nil {
		utils.sdkLog("Unmarshal failed, ", resp, err.Error())
		return
	}

	if userResp.ErrCode != 0 {
		utils.sdkLog("errcode: ", userResp.ErrCode, "errmsg:", userResp.ErrMsg)
		return
	}

	if len(userResp.Data) == 0 {
		utils.sdkLog("failed, no user : ", u.loginUserID)
		return
	}

	err = u.updateFriendInfo(userResp.Data[0].Uid, userResp.Data[0].Name, userResp.Data[0].Icon, userResp.Data[0].Gender, userResp.Data[0].Mobile, userResp.Data[0].Birth, userResp.Data[0].Email, userResp.Data[0].Ex)
	if err != nil {
		utils.sdkLog("  db change failed", err.Error())
		return
	}

	jsonInfo, err := json.Marshal(userResp.Data[0])
	if err != nil {
		utils.sdkLog("  marshal failed", err.Error())
		return
	}

	u.friendListener.OnFriendInfoChanged(string(jsonInfo))
}
